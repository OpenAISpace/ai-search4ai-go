package units

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// SearchResult represents a single search result
type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// SearchResponse represents the response from a search
type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// Search performs a search using the configured search service
func Search(query string) (string, error) {
	fmt.Printf("正在使用查询进行自定义搜索: %s\n", query)

	searchService := os.Getenv("SEARCH_SERVICE")
	if searchService == "" {
		searchService = "duckduckgo" // Default to DuckDuckGo
	}

	var results []SearchResult
	var err error

	switch searchService {
	case "search1api":
		results, err = searchWithSearch1API(query)
	case "google":
		results, err = searchWithGoogle(query)
	case "bing":
		results, err = searchWithBing(query)
	case "serpapi":
		results, err = searchWithSerpAPI(query)
	case "serper":
		results, err = searchWithSerper(query)
	case "duckduckgo":
		results, err = searchWithDuckDuckGo(query)
	case "searxng":
		results, err = searchWithSearXNG(query)
	default:
		return "", fmt.Errorf("不支持的搜索服务: %s", searchService)
	}

	if err != nil {
		return "", fmt.Errorf("搜索失败: %v", err)
	}

	response := SearchResponse{Results: results}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %v", err)
	}

	fmt.Println("自定义搜索服务调用完成")
	return string(jsonData), nil
}

func searchWithSearch1API(query string) ([]SearchResult, error) {
	apiKey := os.Getenv("SEARCH1API_KEY")
	maxResults := os.Getenv("MAX_RESULTS")
	if maxResults == "" {
		maxResults = "10"
	}

	reqBody := map[string]string{
		"query":         query,
		"max_results":   maxResults,
		"crawl_results": "0",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.search1api.com/search/", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func searchWithGoogle(query string) ([]SearchResult, error) {
	cx := os.Getenv("GOOGLE_CX")
	apiKey := os.Getenv("GOOGLE_KEY")
	maxResults := os.Getenv("MAX_RESULTS")
	if maxResults == "" {
		maxResults = "10"
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?cx=%s&key=%s&q=%s",
		url.QueryEscape(cx),
		url.QueryEscape(apiKey),
		url.QueryEscape(query))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleResp struct {
		Items []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range googleResp.Items {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		})
	}

	return results[:min(len(results), parseInt(maxResults))], nil
}

func searchWithBing(query string) ([]SearchResult, error) {
	apiKey := os.Getenv("BING_KEY")
	maxResults := os.Getenv("MAX_RESULTS")
	if maxResults == "" {
		maxResults = "10"
	}

	req, err := http.NewRequest("GET", "https://api.bing.microsoft.com/v7.0/search?q="+url.QueryEscape(query), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bingResp struct {
		WebPages struct {
			Value []struct {
				Name    string `json:"name"`
				URL     string `json:"url"`
				Snippet string `json:"snippet"`
			} `json:"value"`
		} `json:"webPages"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&bingResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range bingResp.WebPages.Value {
		results = append(results, SearchResult{
			Title:   item.Name,
			Link:    item.URL,
			Snippet: item.Snippet,
		})
	}

	return results[:min(len(results), parseInt(maxResults))], nil
}

func searchWithSerpAPI(query string) ([]SearchResult, error) {
	apiKey := os.Getenv("SERPAPI_KEY")
	maxResults := os.Getenv("MAX_RESULTS")
	if maxResults == "" {
		maxResults = "10"
	}

	apiURL := fmt.Sprintf("https://serpapi.com/search?api_key=%s&engine=google&q=%s&google_domain=google.com",
		url.QueryEscape(apiKey),
		url.QueryEscape(query))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var serpResp struct {
		Organic []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"organic_results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&serpResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range serpResp.Organic {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		})
	}

	return results[:min(len(results), parseInt(maxResults))], nil
}

func searchWithSerper(query string) ([]SearchResult, error) {
	apiKey := os.Getenv("SERPER_KEY")
	gl := os.Getenv("GL")
	if gl == "" {
		gl = "us"
	}
	hl := os.Getenv("HL")
	if hl == "" {
		hl = "en"
	}

	reqBody := map[string]string{
		"q":  query,
		"gl": gl,
		"hl": hl,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://google.serper.dev/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var serperResp struct {
		Organic []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"organic"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&serperResp); err != nil {
		return nil, err
	}

	maxResults := parseInt(os.Getenv("MAX_RESULTS"))
	if maxResults == 0 {
		maxResults = 10
	}

	var results []SearchResult
	for _, item := range serperResp.Organic {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		})
	}

	return results[:min(len(results), maxResults)], nil
}

func searchWithDuckDuckGo(query string) ([]SearchResult, error) {
	maxResults := os.Getenv("MAX_RESULTS")
	if maxResults == "" {
		maxResults = "10"
	}

	reqBody := map[string]string{
		"q":           query,
		"max_results": maxResults,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://ddg.search2ai.online/search", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var duckResp struct {
		Results []struct {
			Title string `json:"title"`
			Href  string `json:"href"`
			Body  string `json:"body"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&duckResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range duckResp.Results {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Href,
			Snippet: item.Body,
		})
	}

	return results, nil
}

func searchWithSearXNG(query string) ([]SearchResult, error) {
	baseURL := os.Getenv("SEARXNG_BASE_URL")
	maxResults := parseInt(os.Getenv("MAX_RESULTS"))
	if maxResults == 0 {
		maxResults = 10
	}

	apiURL := fmt.Sprintf("%s/search?q=%s&category=general&format=json",
		baseURL,
		url.QueryEscape(query))

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searxResp struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searxResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range searxResp.Results {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.URL,
			Snippet: item.Content,
		})
	}

	return results[:min(len(results), maxResults)], nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
