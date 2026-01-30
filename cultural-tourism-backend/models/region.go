// File: models/region.go
package models

// Region 区域数据模型
type Region struct {
	ID        string `json:"_id,omitempty"`     // TCB 自动生成的 ID (String)
	OpenID    string `json:"_openid,omitempty"` // 系统字段
	Name      string `json:"name"`              // 区域名称
	Status    int    `json:"status"`            // 1: 启用, 0: 禁用
	Sort      int    `json:"sort"`              // 排序权重 (值越大越靠前)
	CreatedAt string `json:"created_at"`        // 创建时间
	UpdatedAt string `json:"updated_at"`        // 更新时间
}
