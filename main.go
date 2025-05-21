package main

import (
	"go-book/config"
	"go-book/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 确保在main函数结束时关闭数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}
	defer sqlDB.Close()

	// 创建Gin实例
	r := gin.Default()

	// 初始化路由
	routes.SetupRoutes(r, db)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}