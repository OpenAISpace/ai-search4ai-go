package api

type ChatCompletionRequest struct {
	Model      string                   `json:"model"`
	Messages   []map[string]interface{} `json:"messages"`
	MaxTokens  int                      `json:"max_tokens"`
	Tools      []map[string]interface{} `json:"tools,omitempty"`
	ToolChoice string                   `json:"tool_choice,omitempty"`
	Stream     bool                     `json:"stream"`
}

// ChatCompletionResponse represents the response structure from OpenAI
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int                    `json:"index"`
		Message      map[string]interface{} `json:"message"`
		FinishReason string                 `json:"finish_reason"`
	} `json:"choices"`
}

// ToolCall represents a tool call from OpenAI
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// Message represents a chat message
type Message struct {
	Role       string                   `json:"role"`
	Content    string                   `json:"content"`
	ToolCalls  []map[string]interface{} `json:"tool_calls,omitempty"`
	Name       string                   `json:"name,omitempty"`
	ToolCallID string                   `json:"tool_call_id,omitempty"`
}

// Delta represents a streaming delta update
type Delta struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

// StreamChoice represents a choice in the streaming response
type StreamChoice struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}
