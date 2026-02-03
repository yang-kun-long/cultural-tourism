// File: main.go
package main

import (
	"cultural-tourism-backend/routes"
	"cultural-tourism-backend/tcb" // 引入 tcb

	_ "cultural-tourism-backend/docs"

	"github.com/gin-gonic/gin"
)

// @title           数字文旅后端 API
// @version         1.0
// @description     基于 Go + Gin + 腾讯云开发构建的 RESTful API
// @BasePath        /api
// @schemes         https http
func main() {
	// 1. 初始化云开发 HTTP 客户端
	tcb.Init()

	// 2. 初始化 Gin
	r := gin.Default()

	// 3. 注册路由
	routes.RegisterRoutes(r)

	// 4. 启动
	r.Run(":8080")
}
