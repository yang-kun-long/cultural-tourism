// File: config/config.go
package config

import (
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	ServerPort   string
}

var AppConfig = &Config{
	// =======================================================
	// 【请修改】填入腾讯云开发控制台获取的 MongoDB 连接字符串
	// =======================================================
	MongoURI: "mongodb://user:password@address:port/test?authSource=admin&replicaSet=replica",

	DatabaseName: "cultural_tourism_db",
	ServerPort:   getEnv("PORT", "8080"),
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return ":" + value
	}
	return ":" + fallback
}
