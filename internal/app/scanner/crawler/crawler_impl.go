package crawler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gw-gong/key-spy/internal/config/scanner/localcfg"
	"github.com/gw-gong/key-spy/internal/pkg/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/gw-gong/gwkit-go/log"
)

type crawler struct {
	cfg        *localcfg.Config
	httpClient *http.Client
	visited    map[string]bool
	visitedMu  sync.Mutex
	results    []*model.ScanResult
	resultsMu  sync.Mutex
	semaphore  chan struct{}
}

func NewCrawler(cfg *localcfg.Config) Crawler {
	return &crawler{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Scanner.RequestTimeoutMs) * time.Millisecond,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
		visited:   make(map[string]bool),
		results:   make([]*model.ScanResult, 0),
		semaphore: make(chan struct{}, cfg.Scanner.MaxConcurrent),
	}
}

func (c *crawler) Scan(ctx context.Context) (*model.ScanReport, error) {
	startTime := time.Now()

	log.Infoc(ctx, "Starting scan",
		log.Str("target_url", c.cfg.Scanner.TargetURL),
		log.Any("keywords", c.cfg.Scanner.Keywords),
		log.Int("max_depth", c.cfg.Scanner.MaxDepth),
	)

	// 开始爬取
	var wg sync.WaitGroup
	wg.Add(1)
	go c.crawl(ctx, c.cfg.Scanner.TargetURL, 0, &wg)
	wg.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// 构建报告
	c.resultsMu.Lock()
	matchResults := make([]*model.ScanResult, 0)
	for _, r := range c.results {
		if r.TotalCount > 0 {
			matchResults = append(matchResults, r)
		}
	}
	totalPages := len(c.results)
	c.resultsMu.Unlock()

	report := &model.ScanReport{
		TargetURL:  c.cfg.Scanner.TargetURL,
		Keywords:   c.cfg.Scanner.Keywords,
		StartTime:  startTime.Format("2006-01-02 15:04:05"),
		EndTime:    endTime.Format("2006-01-02 15:04:05"),
		Duration:   duration.String(),
		TotalPages: totalPages,
		MatchPages: len(matchResults),
		Results:    matchResults,
	}

	log.Infoc(ctx, "Scan completed",
		log.Int("total_pages", totalPages),
		log.Int("match_pages", len(matchResults)),
		log.Str("duration", duration.String()),
	)

	return report, nil
}

func (c *crawler) crawl(ctx context.Context, pageURL string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 检查深度限制
	if depth > c.cfg.Scanner.MaxDepth {
		return
	}

	// 规范化 URL
	normalizedURL := c.normalizeURL(pageURL)
	if normalizedURL == "" {
		return
	}

	// 检查是否已访问
	c.visitedMu.Lock()
	if c.visited[normalizedURL] {
		c.visitedMu.Unlock()
		return
	}
	c.visited[normalizedURL] = true
	c.visitedMu.Unlock()

	// 检查 URL 是否属于目标域名
	if !c.isSameDomain(normalizedURL) {
		return
	}

	// 获取信号量
	select {
	case c.semaphore <- struct{}{}:
		defer func() { <-c.semaphore }()
	case <-ctx.Done():
		return
	}

	// 请求间隔
	time.Sleep(time.Duration(c.cfg.Scanner.RequestIntervalMs) * time.Millisecond)

	log.Debugc(ctx, "Crawling page", log.Str("url", normalizedURL), log.Int("depth", depth))

	// 获取页面内容
	body, links, err := c.fetchPage(ctx, normalizedURL)
	if err != nil {
		log.Warnc(ctx, "Failed to fetch page", log.Str("url", normalizedURL), log.Err(err))
		c.addResult(&model.ScanResult{
			URL:   normalizedURL,
			Depth: depth,
			Error: err.Error(),
		})
		return
	}

	// 搜索关键词
	result := c.searchKeywords(normalizedURL, body, depth)
	c.addResult(result)

	if result.TotalCount > 0 {
		log.Infoc(ctx, "Found keywords",
			log.Str("url", normalizedURL),
			log.Any("keywords", result.Keywords),
			log.Int("total_count", result.TotalCount),
		)
	}

	// 递归爬取链接
	for _, link := range links {
		absoluteURL := c.resolveURL(normalizedURL, link)
		if absoluteURL != "" {
			wg.Add(1)
			go c.crawl(ctx, absoluteURL, depth+1, wg)
		}
	}
}

func (c *crawler) fetchPage(ctx context.Context, pageURL string) (body string, links []string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Set("User-Agent", c.cfg.Scanner.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// 只处理 HTML 内容
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") && !strings.Contains(contentType, "application/xhtml") {
		return "", nil, nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	body = string(bodyBytes)

	// 解析链接
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return body, nil, nil
	}

	links = make([]string, 0)
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, href)
		}
	})

	return body, links, nil
}

func (c *crawler) searchKeywords(pageURL, body string, depth int) *model.ScanResult {
	result := &model.ScanResult{
		URL:           pageURL,
		KeywordCounts: make(map[string]int),
		Keywords:      make([]string, 0),
		Depth:         depth,
	}

	for _, keyword := range c.cfg.Scanner.Keywords {
		// 使用正则表达式进行不区分大小写的搜索
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
		matches := re.FindAllString(body, -1)
		count := len(matches)

		if count > 0 {
			result.KeywordCounts[keyword] = count
			result.Keywords = append(result.Keywords, keyword)
			result.TotalCount += count
		}
	}

	return result
}

func (c *crawler) addResult(result *model.ScanResult) {
	c.resultsMu.Lock()
	c.results = append(c.results, result)
	c.resultsMu.Unlock()
}

func (c *crawler) normalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	// 移除 fragment
	parsed.Fragment = ""

	// 确保有 scheme
	if parsed.Scheme == "" {
		parsed.Scheme = "https"
	}

	return parsed.String()
}

func (c *crawler) isSameDomain(pageURL string) bool {
	targetParsed, err := url.Parse(c.cfg.Scanner.TargetURL)
	if err != nil {
		return false
	}

	pageParsed, err := url.Parse(pageURL)
	if err != nil {
		return false
	}

	return targetParsed.Host == pageParsed.Host
}

func (c *crawler) resolveURL(baseURL, href string) string {
	if href == "" {
		return ""
	}

	// 跳过非 HTTP 链接
	if strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") ||
		strings.HasPrefix(href, "#") {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	resolved := base.ResolveReference(ref)
	resolved.Fragment = ""

	return resolved.String()
}
