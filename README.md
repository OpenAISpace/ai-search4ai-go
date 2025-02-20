# Search4AI-Go

Search4AI-Go 是一个用 Go 语言实现的 AI 搜索代理服务，它允许大型语言模型通过 API 调用来获取实时互联网信息。这是 [search2ai](https://github.com/liyown/search2ai) 项目的 Go 语言版本。

## 功能特点

- 支持多个搜索引擎服务：
  - DuckDuckGo（默认，无需 API 密钥）
  - Google Custom Search
  - Bing Web Search
  - SerpAPI
  - Serper
  - Search1API
  - SearXNG（自托管选项）
- 网页内容抓取和分析
- 支持流式响应
- 完整的 CORS 支持
- 与 OpenAI API 格式兼容
- 搜索结果实时返回到流式响应中

## 快速开始

### 安装

1. 克隆仓库：

```bash
git clone https://github.com/liyown/search4ai-go.git
cd search4ai-go
```

2. 安装依赖：

```bash
go mod download
```

3. 配置环境变量：

复制 `.env.example` 文件为 `.env` 并根据需要修改配置：

```bash
cp .env.example .env
```

### 配置

在 `.env` 文件中设置以下配置项：

```env
# 服务器配置
PORT=3014                          # 服务器端口
APIBASE=https://api.openai.com     # AI 模型 API 基础 URL

# 搜索配置
SEARCH_SERVICE=duckduckgo         # 默认搜索服务
MAX_RESULTS=10                    # 每次搜索返回的最大结果数

# Google 搜索配置（如果使用 Google）
GOOGLE_CX=your_google_cx          # Google 自定义搜索引擎 ID
GOOGLE_KEY=your_google_api_key    # Google API 密钥

# 其他搜索服务配置（根据需要取消注释）
#BING_KEY=your_bing_api_key       # Bing API 密钥
#SERPAPI_KEY=your_serpapi_key     # SerpAPI 密钥
#SERPER_KEY=your_serper_key       # Serper API 密钥
#SEARCH1API_KEY=your_search1api_key # Search1API 密钥
#SEARXNG_BASE_URL=your_searxng_url # SearXNG 自托管 URL
```

### 运行

启动服务器：

```bash
go run main.go
```

服务器将在配置的端口上运行（默认为 3014）。

## API 使用

### 1. 基本搜索

发送 POST 请求到 `/v1/chat/completions` 端点：

```json
{
  "model": "moonshot-v1-128k",
  "messages": [
    {
      "role": "system",
      "content": "你是一个有用的助手。当用户请求实时信息（例如日期、天气或新闻）时，使用函数调用来检索相关数据。如有必要，递归调用函数来收集其他详细信息。"
    },
    {
      "role": "user",
      "content": "最近的世界新闻有哪些？"
    }
  ],
  "stream": true,
  "enabledTools": {
    "search": true
  }
}
```

响应示例（流式响应的一个数据块）：
```json
{
    "id": "chatcmpl-123",
    "object": "chat.completion.chunk",
    "created": 1677652288,
    "model": "moonshot-v1-128k",
    "choices": [{
        "index": 0,
        "delta": {
            "content": "根据最新的新闻报道，"
        },
        "finish_reason": null
    }],
    "system_fingerprint": "fp-123",
    "search_results": [{
        "title": "Latest World News - Reuters",
        "link": "https://www.reuters.com/world/",
        "snippet": "Get the latest world news coverage..."
    }]
}
```

### 2. 网页抓取

使用网页抓取功能：

```json
{
  "model": "moonshot-v1-128k",
  "messages": [
    {
      "role": "system",
      "content": "你是一个有用的助手。需要详细信息时，使用爬虫功能获取网页内容。"
    },
    {
      "role": "user",
      "content": "帮我获取并总结这个网页的内容：https://example.com"
    }
  ],
  "stream": true,
  "enabledTools": {
    "crawler": true
  }
}
```

### 工具说明

1. **search 工具**
   - 用于获取实时互联网信息
   - 自动在对话中使用，无需手动指定参数
   - 搜索结果会在流式响应中实时返回

2. **crawler 工具**
   - 用于抓取和分析特定网页内容
   - 自动在对话中使用，无需手动指定参数
   - 支持大多数常见网页格式

### 使用提示

1. **系统提示（System Prompt）**
   - 建议在 system 消息中说明助手的行为，特别是何时使用搜索或爬虫功能
   - 例如："当需要实时信息时使用搜索功能"或"需要详细内容时使用爬虫功能"

2. **工具启用**
   - 使用 `enabledTools` 字段控制可用的工具
   - 可以同时启用多个工具：`{"search": true, "crawler": true}`

3. **流式响应**
   - 设置 `stream: true` 获取实时响应
   - 搜索结果会在 `search_results` 字段中返回
   - 每个数据块都包含完整的元数据

## 搜索服务说明

1. **DuckDuckGo**（默认）
   - 无需 API 密钥
   - 适合一般用途

2. **Google Custom Search**
   - 需要 Google Custom Search Engine ID 和 API 密钥
   - 提供高质量的搜索结果
   - 每日请求限制

3. **Bing Web Search**
   - 需要 Bing API 密钥
   - 提供全面的搜索结果

4. **SerpAPI**
   - 需要 SerpAPI 密钥
   - 提供多个搜索引擎的结果

5. **Serper**
   - 需要 Serper API 密钥
   - Google 搜索结果的替代方案

6. **Search1API**
   - 需要 Search1API 密钥
   - 提供自定义搜索功能

7. **SearXNG**
   - 自托管选项
   - 完全可控的搜索引擎元搜索引擎

## 注意事项

1. 流式响应中的搜索结果会实时返回
2. 每个搜索服务可能有不同的速率限制和定价
3. 建议在生产环境中使用环境变量管理 API 密钥
4. 确保您的 API 密钥有足够的配额

## 贡献

欢迎提交 Pull Requests 和 Issues！

## 许可证

MIT License 