package stream

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	ID                string                   `json:"id"`
	Object            string                   `json:"object"`
	Created           int64                    `json:"created"`
	Model             string                   `json:"model"`
	Choices           []StreamChoice           `json:"choices"`
	SystemFingerprint string                   `json:"system_fingerprint"`
	SearchResults     []map[string]interface{} `json:"search_results,omitempty"`
}

// StreamChoice represents a choice in the streaming response
type StreamChoice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

// Delta represents a streaming delta update
type Delta struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool call from the model
type ToolCall struct {
	Index    int      `json:"index"`
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents the function details in a tool call
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Message represents a chat message
type Message struct {
	Role       string                   `json:"role"`
	Content    string                   `json:"content"`
	ToolCalls  []map[string]interface{} `json:"tool_calls,omitempty"`
	Name       string                   `json:"name,omitempty"`
	ToolCallID string                   `json:"tool_call_id,omitempty"`
}
