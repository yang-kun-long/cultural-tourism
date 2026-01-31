// File: models/product.go
package models

// Product 商品导流模型
// 对应 PRD 5.5 (商品导流) & 7.1 (数据对象)
type Product struct {
	ID        string  `json:"_id,omitempty"`     // TCB 自动生成的 ID
	OpenID    string  `json:"_openid,omitempty"` // 系统字段：创建者
	Name      string  `json:"name"`              // 商品名称
	Image     string  `json:"image"`             // 商品图片 URL
	Price     float64 `json:"price"`             // 商品价格 (仅展示，无支付)
	JumpAppID string  `json:"jump_app_id"`       // 跳转小程序 AppID
	JumpPath  string  `json:"jump_path"`         // 跳转路径
	CreatedAt string  `json:"created_at"`        // 业务创建时间
	UpdatedAt string  `json:"updated_at"`        // 业务更新时间
}

// ProductQuery 商品列表筛选参数
type ProductQuery struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=10"`
}
