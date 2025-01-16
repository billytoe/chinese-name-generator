package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/billytoe/chinese-name-generator/config"
)

// OllamaClient 处理与Ollama服务的通信
type OllamaClient struct {
	config *config.OllamaConfig
	client *http.Client
}

// OllamaRequest 表示发送给Ollama的请求
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // 设置为false以禁用流式输出
	Format   string    `json:"format"` // 指定输出格式
}

// Message 表示对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse 表示Ollama的响应
type OllamaResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

// NewOllamaClient 创建新的Ollama客户端
func NewOllamaClient(config *config.OllamaConfig) *OllamaClient {
	return &OllamaClient{
		config: config,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// GenerateChineseNames 生成中文名字建议
func (c *OllamaClient) GenerateChineseNames(englishName string) ([]NameSuggestion, error) {
	prompt := fmt.Sprintf(`你是一个专业的中文起名专家。请为英文名"%s"推荐10个富有文化内涵的中文名。
请严格按照以下JSON格式返回，不要包含任何其他内容：
{
  "suggestions": [
    {
      "chinese_name": "中文名",
      "pinyin": "拼音",
      "meaning": "含义解释",
      "english_explanation": "英文说明",
      "cultural_context": "文化内涵"
    }
  ]
}

要求：
1. 音韵和谐，尽可能与英文发音相近
2. 符合中国传统起名规范
3. 字义优美，组合得当
4. 避免文化禁忌`, englishName)

	reqBody := OllamaRequest{
		Model: c.config.EmbeddingModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,  // 禁用流式输出
		Format: "json", // 指定JSON格式输出
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	resp, err := c.client.Post(
		c.config.URL+"/api/chat",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %v", err)
	}

	// 记录原始响应以便调试
	log.Printf("Ollama raw response: %s", string(respBody))

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, fmt.Errorf("decode response failed: %v", err)
	}

	// 从消息内容中提取JSON
	jsonContent := extractJSON(ollamaResp.Message.Content)
	if jsonContent == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	// 解析JSON响应
	var result struct {
		Suggestions []NameSuggestion `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("parse suggestions failed: %v, JSON: %s", err, jsonContent)
	}

	return result.Suggestions, nil
}

// extractJSON 从文本中提取JSON部分
func extractJSON(text string) string {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end < start {
		return ""
	}
	return text[start : end+1]
}

// NameSuggestion 表示一个中文名字建议
type NameSuggestion struct {
	ChineseName        string `json:"chinese_name"`
	Pinyin             string `json:"pinyin"`
	Meaning            string `json:"meaning"`
	EnglishExplanation string `json:"english_explanation"`
	CulturalContext    string `json:"cultural_context"`
}
