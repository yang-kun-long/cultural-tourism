// File: models/theme.go
package models

// Theme 旅拍主题模型
// 对应 PRD 5.1 (主题展示) & 7.1 (数据对象)
type Theme struct {
	ID        string `json:"_id,omitempty"`     // TCB 自动生成的 ID
	OpenID    string `json:"_openid,omitempty"` // 系统字段：创建者
	Name      string `json:"name"`              // 主题名称 (如: "汉服打卡")
	Cover     string `json:"cover"`             // 封面图 URL
	Desc      string `json:"desc"`              // 主题简介
	RegionID  string `json:"region_id"`         // 所属区域 ID (用于实现 PRD 5.1.1 区域优先推荐)
	Sort      int    `json:"sort"`              // 排序权重 (默认 100)
	Status    int    `json:"status"`            // 1:启用 0:禁用
	CreatedAt string `json:"created_at"`        // 业务创建时间
	UpdatedAt string `json:"updated_at"`        // 业务更新时间
}

// ThemeQuery 主题列表筛选参数
type ThemeQuery struct {
	RegionID string `form:"region_id"`        // 核心筛选：按区域
	Status   int    `form:"status,default=1"` // 默认只查启用
	Page     int    `form:"page,default=1"`
	Size     int    `form:"size,default=10"`
}
