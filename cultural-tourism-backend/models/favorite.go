package models

// Favorite 收藏模型
type Favorite struct {
	ID           string `json:"_id,omitempty" bson:"_id,omitempty"`
	OpenID       string `json:"_openid,omitempty" bson:"_openid,omitempty"` // 系统字段：用户标识
	ResourceType string `json:"resource_type" bson:"resource_type" binding:"required"` // 资源类型: theme/poi/product
	ResourceID   string `json:"resource_id" bson:"resource_id" binding:"required"`     // 资源ID
	CreatedAt    string `json:"created_at" bson:"created_at"`
	UpdatedAt    string `json:"updated_at" bson:"updated_at"`
}

// FavoriteCreateRequest 收藏创建请求
type FavoriteCreateRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=theme poi product"` // 枚举验证
	ResourceID   string `json:"resource_id" binding:"required"`
}

// FavoriteListRequest 收藏列表请求
type FavoriteListRequest struct {
	ResourceType string `form:"resource_type" binding:"omitempty,oneof=theme poi product"` // 可选筛选
	Page         int    `form:"page" binding:"omitempty,min=1"`
	Size         int    `form:"size" binding:"omitempty,min=1,max=100"`
}
