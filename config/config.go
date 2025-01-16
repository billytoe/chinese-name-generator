package config

// OllamaConfig 存储Ollama的配置信息
type OllamaConfig struct {
	URL            string `json:"url"`
	EmbeddingModel string `json:"embedding_model"`
}

// Config 应用程序配置
type Config struct {
	Ollama OllamaConfig `json:"ollama"`
}

// GetConfig 返回默认配置
func GetConfig() *Config {
	return &Config{
		Ollama: OllamaConfig{
			URL:            "http://210.73.217.30:28091",
			EmbeddingModel: "glm4:9b",
		},
	}
}
