package model

// ScanResult 表示单个页面的扫描结果
type ScanResult struct {
	URL           string            `json:"url"`            // 页面 URL
	KeywordCounts map[string]int    `json:"keyword_counts"` // 关键词出现次数
	TotalCount    int               `json:"total_count"`    // 总出现次数
	Keywords      []string          `json:"keywords"`       // 出现的关键词列表
	Depth         int               `json:"depth"`          // 页面深度
	Error         string            `json:"error,omitempty"` // 错误信息（如有）
}

// ScanReport 表示完整的扫描报告
type ScanReport struct {
	TargetURL   string        `json:"target_url"`   // 目标网站
	Keywords    []string      `json:"keywords"`     // 搜索的关键词列表
	StartTime   string        `json:"start_time"`   // 开始时间
	EndTime     string        `json:"end_time"`     // 结束时间
	Duration    string        `json:"duration"`     // 耗时
	TotalPages  int           `json:"total_pages"`  // 扫描的总页面数
	MatchPages  int           `json:"match_pages"`  // 匹配的页面数
	Results     []*ScanResult `json:"results"`      // 匹配的结果
	ErrorCount  int           `json:"error_count"`  // 错误数
}
