package stream

// ToolCallCollector collects and manages tool calls
type ToolCallCollector struct {
	currentToolCall map[string]interface{}
	toolCalls       []map[string]interface{}
	toolCallResults []map[string]interface{}
}

// NewToolCallCollector creates a new tool call collector
func NewToolCallCollector() *ToolCallCollector {
	return &ToolCallCollector{
		currentToolCall: make(map[string]interface{}),
		toolCalls:       make([]map[string]interface{}, 0),
	}
}

// CollectToolCall collects a tool call
func (tc *ToolCallCollector) CollectToolCall(call ToolCall) {
	if tc.currentToolCall["id"] == nil {
		tc.currentToolCall = map[string]interface{}{
			"id":   call.ID,
			"type": call.Type,
			"function": map[string]interface{}{
				"name":      call.Function.Name,
				"arguments": call.Function.Arguments,
			},
		}
	} else {
		if call.ID != "" {
			tc.currentToolCall["id"] = call.ID
		}
		if call.Function.Name != "" {
			tc.currentToolCall["function"].(map[string]interface{})["name"] = call.Function.Name
		}
		if call.Function.Arguments != "" {
			args := tc.currentToolCall["function"].(map[string]interface{})["arguments"].(string)
			args += call.Function.Arguments
			tc.currentToolCall["function"].(map[string]interface{})["arguments"] = args
		}
	}
}

// GetToolCalls returns the collected tool calls
func (tc *ToolCallCollector) GetToolCalls() []map[string]interface{} {
	if tc.currentToolCall["id"] != nil {
		tc.toolCalls = append(tc.toolCalls, tc.currentToolCall)
		tc.currentToolCall = make(map[string]interface{})
	}
	return tc.toolCalls
}

// GetToolCallResults returns the collected tool call results
func (tc *ToolCallCollector) GetToolCallResults() []map[string]interface{} {
	return tc.toolCallResults
}
