package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/liyown/search4ai-go/units"
)

func executeToolCalls(toolCalls []interface{}) ([]map[string]interface{}, error) {
	var toolResults []map[string]interface{}
	for _, tc := range toolCalls {
		toolCall, ok := tc.(map[string]interface{})
		if !ok {
			continue
		}

		result, err := executeToolCall(toolCall)
		if err != nil {
			log.Printf("Error executing tool call: %v", err)
			continue
		}

		toolResults = append(toolResults, map[string]interface{}{
			"tool_call_id": toolCall["id"],
			"role":         "tool",
			"name":         toolCall["function"].(map[string]interface{})["name"],
			"content":      result,
		})
	}
	return toolResults, nil
}

// executeToolCall executes a tool call and returns the result
func executeToolCall(toolCall map[string]interface{}) (string, error) {
	function, ok := toolCall["function"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid tool call format")
	}

	name, ok := function["name"].(string)
	if !ok {
		return "", fmt.Errorf("invalid function name")
	}

	arguments, ok := function["arguments"].(string)
	if !ok {
		return "", fmt.Errorf("invalid function arguments")
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("error parsing arguments: %v", err)
	}

	switch name {
	case "search":
		query, ok := args["query"].(string)
		if !ok {
			return "", fmt.Errorf("invalid search query")
		}
		return units.Search(query)

	case "crawler":
		url, ok := args["url"].(string)
		if !ok {
			return "", fmt.Errorf("invalid crawler url")
		}
		return units.Crawler(url)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// buildTools creates the tools configuration
func buildTools(enabledTools map[string]bool) []map[string]interface{} {
	tools := []map[string]interface{}{
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "search",
				"description": "搜索互联网获取实时信息。当你需要查找当前信息时使用此功能，例如日期、天气、新闻，或者可能不在你训练数据中的事实。搜索结果将包含相关网页的标题、链接和摘要。",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "搜索查询。应该具体且聚焦于所需信息。使用可能出现在相关结果中的关键词和短语。",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "crawler",
				"description": "提取和分析特定网页URL的内容。当你需要从特定网页获取详细信息时使用此功能，包括其文本内容、标题和元数据。这对于在通过搜索找到相关URL后获取深入信息很有用。",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "要分析的网页的完整URL。必须是以http://或https://开头的有效、可访问的网址。",
						},
					},
					"required": []string{"url"},
				},
			},
		},
	}

	if enabledTools == nil {
		return tools
	}

	var filteredTools []map[string]interface{}
	for _, tool := range tools {
		name := tool["function"].(map[string]interface{})["name"].(string)
		if enabled, exists := enabledTools[name]; !exists || enabled {
			filteredTools = append(filteredTools, tool)
		}
	}
	return filteredTools
}
