# Key-Spy

网站关键词扫描工具，支持定时扫描目标网站的所有页面，查找指定关键词并生成报告。

## 配置

编辑 `config/scanner/localcfg/live.yaml`：

```yaml
scanner:
  target_url: "https://example.com"   # 目标网站
  keywords:                           # 关键词列表
    - "关键词1"
    - "关键词2"
  max_depth: 5                        # 最大爬取深度
  request_interval_ms: 1000           # 请求间隔（毫秒）

cron:
  spec: "0 2 * * *"                   # cron 表达式
  enabled: true                       # true 启用定时任务，false 单次执行

output:
  dir: "/data/key-spy/output"         # 输出目录
```

## 部署

```bash
cd deploy/scanner

# 构建镜像
make docker

# 启动服务
make docker-up

# 停止服务
make docker-down

# 查看所有命令
make help
```

## 输出

扫描报告保存在 `output/` 目录，格式为 `scan_result_20260119_150000.txt`。
