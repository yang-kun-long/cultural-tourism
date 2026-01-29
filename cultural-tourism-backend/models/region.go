// File: models/region.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Region 区域数据模型
type Region struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name   string             `bson:"name" json:"name" binding:"required"` // 区域名称 (如: "西湖区", "主景区")
	Status int                `bson:"status" json:"status"`                // 1: 启用, 0: 禁用
	Sort   int                `bson:"sort" json:"sort"`                    // 排序权重
}
