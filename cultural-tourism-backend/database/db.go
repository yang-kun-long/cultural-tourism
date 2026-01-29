// File: database/db.go
package database

import (
	"context"
	"log"
	"time"

	"cultural-tourism-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ 数据库客户端创建失败: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	DB = client.Database(config.AppConfig.DatabaseName)
	log.Println("✅ 数据库连接成功")
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
