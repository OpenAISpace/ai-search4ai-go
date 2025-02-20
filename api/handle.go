package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liyown/search4ai-go/stream"
)

func handleStreamingResponse(c *gin.Context, resp *http.Response, req *ChatCompletionRequest) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(resp.StatusCode)
	var searchResults []map[string]interface{}

	for {
		processor := stream.NewProcessor(c.Writer)
		message, collectedTools, needsToolExecution := processor.ProcessStream(resp.Body, searchResults)

		// If no tool execution is needed, we're done
		if !needsToolExecution {
			fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
			return
		}

		// Add the assistant's tool calls message to the conversation
		req.Messages = append(req.Messages, map[string]interface{}{
			"role":       "assistant",
			"tool_calls": message.ToolCalls,
		})

		// Execute collected tool calls
		toolCallsInterface := make([]interface{}, len(collectedTools))
		for i, v := range collectedTools {
			toolCallsInterface[i] = v
		}
		toolResults, err := executeToolCalls(toolCallsInterface)
		if err != nil {
			log.Printf("Error executing tool calls: %v", err)
			fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
			return
		}

		searchResults = toolResults

		// Add tool results to the conversation
		req.Messages = append(req.Messages, toolResults...)

		// Make a new request with the updated context
		newResp, err := forwardToOpenAI(req, strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
		if err != nil {
			log.Printf("Error making recursive request: %v", err)
			fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
			return
		}
		resp.Body.Close()
		resp = newResp

	}
}

func handleNonStreamingResponse(c *gin.Context, resp *http.Response, req *ChatCompletionRequest) {
	var openaiResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing OpenAI response"})
		return
	}

	// Check for tool calls
	if len(openaiResp.Choices) > 0 && openaiResp.Choices[0].Message != nil {
		if toolCalls, ok := openaiResp.Choices[0].Message["tool_calls"].([]interface{}); ok {
			toolResults, err := executeToolCalls(toolCalls)
			if err != nil {
				log.Printf("error executing tool calls: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error executing tool calls"})
				return
			}

			if len(toolResults) > 0 {
				// Add results to messages and make recursive request
				req.Messages = append(req.Messages, openaiResp.Choices[0].Message)
				req.Messages = append(req.Messages, toolResults...)

				body, err := json.Marshal(req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "error preparing tool results"})
					return
				}

				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				handleChatCompletions(c)
				return
			}
		}
	}

	c.JSON(resp.StatusCode, openaiResp)
}

// handleChatCompletions handles the chat completions endpoint
func handleChatCompletions(c *gin.Context) {
	// Validate and prepare request
	req, apiKey, err := prepareRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward request to OpenAI
	resp, err := forwardToOpenAI(req, apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Handle response based on streaming flag
	if req.Stream {
		handleStreamingResponse(c, resp, req)
	} else {
		handleNonStreamingResponse(c, resp, req)
	}
}
