package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupCORS adds CORS middleware to the Gin engine
func setupCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func prepareRequest(c *gin.Context) (*ChatCompletionRequest, string, error) {
	// Get API base URL
	apiBase := os.Getenv("APIBASE")
	if apiBase == "" {
		apiBase = "https://api.openai.com"
	}

	// Get API key from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, "", fmt.Errorf("authorization header is missing")
	}
	apiKey := strings.TrimPrefix(authHeader, "Bearer ")

	// Read and parse request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading request body: %v", err)
	}

	var req ChatCompletionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, "", fmt.Errorf("error parsing request body: %v", err)
	}

	// Add tools if not present
	if req.Tools == nil {
		req.Tools = buildTools(nil)
	}

	return &req, apiKey, nil
}

func forwardToOpenAI(req *ChatCompletionRequest, apiKey string) (*http.Response, error) {
	apiBase := os.Getenv("APIBASE")
	if apiBase == "" {
		apiBase = "https://api.openai.com"
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error preparing request: %v", err)
	}

	client := &http.Client{}
	openaiReq, err := http.NewRequest("POST", apiBase+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenAI request: %v", err)
	}

	openaiReq.Header.Set("Content-Type", "application/json")
	openaiReq.Header.Set("Authorization", "Bearer "+apiKey)

	return client.Do(openaiReq)
}

// StartServer initializes and starts the HTTP server
func StartServer() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin engine
	r := gin.New()

	// Use middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(setupCORS())

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, "<html><head><meta charset=\"UTF-8\"></head><body><h1>欢迎体验search4ai，让你的大模型自由联网！</h1></body></html>")
	})

	// Chat completions endpoint
	r.POST("/v1/chat/completions", handleChatCompletions)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3014"
	}

	// Start server
	log.Printf("Server is listening on port %s", port)
	return r.Run(":" + port)
}
