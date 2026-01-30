// File: models/photo.go
package models

// Photo 用户旅拍照片
// 对应 PRD 5.1.3 (上传) & 7.1 (数据对象)
type Photo struct {
	ID        string `json:"_id,omitempty"`     // TCB 自动生成的 ID
	OpenID    string `json:"_openid,omitempty"` // 系统字段：发布者 (用户ID)
	ThemeID   string `json:"theme_id"`          // 关联的主题 ID
	ImageURL  string `json:"image_url"`         // 照片云存储 URL
	Status    int    `json:"status"`            // 审核状态: 0=待审(默认), 1=通过, 2=拒绝
	LikeCount int    `json:"like_count"`        // 点赞数 (用于热门排序)
	CreatedAt string `json:"created_at"`        // 业务创建时间
	UpdatedAt string `json:"updated_at"`        // 业务更新时间 (审核时间)
}

// PhotoQuery 照片列表筛选参数
type PhotoQuery struct {
	ThemeID string `form:"theme_id"`         // 场景：查看某主题下的瀑布流
	Status  int    `form:"status,default=1"` // 场景：前端默认只展示已过审(1)的照片
	Page    int    `form:"page,default=1"`
	Size    int    `form:"size,default=20"` // 瀑布流通常每页多一点
}
