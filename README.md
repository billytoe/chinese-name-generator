# 智能中文姓名推荐系统

这是一个基于Go语言开发的智能化英文转中文姓名推荐系统，旨在为外国友人提供富有文化内涵和个性化的中文名字选择。

## 功能特点

1. **智能名字生成**
   - 支持英文名/姓名转换为中文名
   - 每次生成10个独特的中文名字方案
   - 确保音韵和谐，尽可能与英文发音相近
   - 符合中国传统起名规范

2. **文化解读**
   - 提供单字解释
   - 说明整体含义
   - 阐述文化内涵
   - 分析个性特征
   - 配备英文说明

## 技术架构

- 后端：Go语言
- AI服务：Ollama API (glm4:9b模型)
- 前端：HTML模板 + Tailwind CSS
- 默认端口：7501

## 快速开始

1. 确保已安装Go 1.16或更高版本
2. 克隆项目到本地
3. 启动服务
   ```bash
   go run main.go
   ```
4. 访问 http://localhost:7501 使用服务

## API文档

### POST /api/generate-names
生成中文名字推荐

请求参数：
```json
{
    "english_name": "string"  // 英文名
}
```

响应示例：
```json
{
    "suggestions": [
        {
            "chinese_name": "米凯乐",
            "pinyin": "Mi Kai Le",
            "meaning": "凯旋欢乐",
            "english_explanation": "One who brings joy and triumph",
            "cultural_context": "象征积极向上，充满活力"
        }
    ]
}
```

## 项目结构

```
.
├── README.md
├── main.go                 // 主程序入口
├── config/                 // 配置文件
│   └── config.go          // 应用配置
├── services/              // 业务逻辑
│   └── ollama_client.go   // Ollama API客户端
├── templates/             // HTML模板
│   └── index.html        // 主页模板
└── static/                // 静态资源文件
```

## 注意事项

- 确保Ollama服务可访问（默认地址：http://210.73.217.30:28091）
- 推荐使用Chrome或Firefox最新版本访问
- 系统响应时间可能因AI服务而有所波动

## 开发计划

- [ ] 支持更多语言的名字转换
- [ ] 添加用户名字收藏功能
- [ ] 优化名字生成算法
- [ ] 增加更多文化解读内容

## 贡献指南

欢迎提交Issue和Pull Request来帮助改进项目。
