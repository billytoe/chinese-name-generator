package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/billytoe/chinese-name-generator/config"
	"github.com/billytoe/chinese-name-generator/services"
)

// NameRequest 表示接收到的名字生成请求
type NameRequest struct {
	EnglishName string `json:"english_name"`
}

var ollamaClient *services.OllamaClient

func init() {
	// 初始化Ollama客户端
	cfg := config.GetConfig()
	ollamaClient = services.NewOllamaClient(&cfg.Ollama)
}

// 处理主页请求
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// 处理名字生成API请求
func handleGenerateNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req NameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}

	suggestions, err := ollamaClient.GenerateChineseNames(req.EnglishName)
	if err != nil {
		log.Printf("生成名字失败: %v", err)
		http.Error(w, "生成名字时发生错误", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
	})
}

func main() {
	// 创建必要的目录
	os.MkdirAll("templates", 0755)
	os.MkdirAll("static", 0755)

	// 设置路由
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/generate-names", handleGenerateNames)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 启动服务器
	log.Printf("服务器启动在 http://localhost:7501")
	if err := http.ListenAndServe(":7501", nil); err != nil {
		log.Fatal(err)
	}
}
