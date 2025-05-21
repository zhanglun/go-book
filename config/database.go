package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() (*gorm.DB, error) {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("加载.env文件失败: %v", err)
	}

	// 从环境变量获取数据库配置信息
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("无效的端口号: %v", err)
	}
	username := os.Getenv("DB_USERNAME")

	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=require TimeZone=Asia/Shanghai",
		host,
		username,
		password,
		port,
	)
	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	return db, nil
}