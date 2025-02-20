package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Processor handles stream processing
type Processor struct {
	writer            io.Writer
	message           *Message
	toolCallCollector *ToolCallCollector
}

// NewProcessor creates a new stream processor
func NewProcessor(w io.Writer) *Processor {
	return &Processor{
		writer: w,
		message: &Message{
			Role:      "assistant",
			Content:   "",
			ToolCalls: make([]map[string]interface{}, 0),
		},
		toolCallCollector: NewToolCallCollector(),
	}
}

// ProcessStream processes the stream and returns the message, collected tool calls, and whether tool execution is needed
func (p *Processor) ProcessStream(body io.ReadCloser, searchResults []map[string]interface{}) (*Message, []map[string]interface{}, bool) {
	scanner := bufio.NewScanner(body)
	chunkCount := 0
	isToolCallMessage := false

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var response StreamResponse
		if err := json.Unmarshal([]byte(data), &response); err != nil {
			log.Printf("Error parsing stream data: %v", err)
			continue
		}

		if len(response.Choices) == 0 {
			continue
		}

		// Check for tool calls in the second chunk
		chunkCount++
		if chunkCount == 2 {
			isToolCallMessage = len(response.Choices[0].Delta.ToolCalls) > 0
		}

		choice := response.Choices[0]
		delta := choice.Delta

		// Only process after we know if it's a tool call message (after second chunk)
		if chunkCount >= 2 {
			if isToolCallMessage {
				p.handleFunctionCall(delta)
			} else {
				p.handleContent(delta, response, searchResults)
			}
		}

		// Check if we're done with this stream
		if choice.FinishReason != "" {
			if choice.FinishReason == "tool_calls" {
				// Return collected tool calls for execution
				return p.message, p.toolCallCollector.GetToolCalls(), true
			}
			// If finish reason is not tool_calls, we're done
			return p.message, nil, false
		}
	}

	return p.message, nil, false
}

func (p *Processor) handleContent(delta Delta, response StreamResponse, searchResults []map[string]interface{}) {
	if delta.Role != "" {
		p.message.Role = delta.Role
	}

	if delta.Content != "" {
		p.message.Content += delta.Content
		// Stream content to client with proper escaping
		escapedContent := strings.Replace(delta.Content, "\"", "\\\"", -1)
		escapedContent = strings.Replace(escapedContent, "\n", "\\n", -1)

		// Include metadata from original response and add search results
		streamResp := StreamResponse{
			ID:                response.ID,
			Object:            response.Object,
			Created:           response.Created,
			Model:             response.Model,
			SystemFingerprint: response.SystemFingerprint,
			Choices: []StreamChoice{
				{
					Delta: Delta{
						Content: escapedContent,
					},
					FinishReason: response.Choices[0].FinishReason,
				},
			},
			SearchResults: searchResults,
		}

		respBytes, err := json.Marshal(streamResp)
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			return
		}

		fmt.Fprintf(p.writer, "data: %s\n\n", string(respBytes))
		if f, ok := p.writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (p *Processor) handleFunctionCall(delta Delta) {
	if len(delta.ToolCalls) > 0 {
		p.toolCallCollector.CollectToolCall(delta.ToolCalls[0])

		if len(p.message.ToolCalls) == 0 {
			p.message.ToolCalls = append(p.message.ToolCalls, map[string]interface{}{
				"id":   delta.ToolCalls[0].ID,
				"type": "function",
				"function": map[string]interface{}{
					"name":      delta.ToolCalls[0].Function.Name,
					"arguments": delta.ToolCalls[0].Function.Arguments,
				},
			})
		} else {
			// Append to existing tool call
			currentTool := p.message.ToolCalls[0]
			if delta.ToolCalls[0].ID != "" {
				currentTool["id"] = delta.ToolCalls[0].ID
			}
			if delta.ToolCalls[0].Function.Name != "" {
				currentTool["function"].(map[string]interface{})["name"] = delta.ToolCalls[0].Function.Name
			}
			if delta.ToolCalls[0].Function.Arguments != "" {
				args := currentTool["function"].(map[string]interface{})["arguments"].(string)
				args += delta.ToolCalls[0].Function.Arguments
				currentTool["function"].(map[string]interface{})["arguments"] = args
			}
		}
	}
}
