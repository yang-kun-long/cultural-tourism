// File: models/comment.go
package models

// Comment 评论/互动
// 对应 PRD 5.4 (评论互动) & 7.1 (数据对象)
type Comment struct {
	ID        string `json:"_id,omitempty"`     // TCB 自动 ID
	OpenID    string `json:"_openid,omitempty"` // 发布者 (系统字段)
	POIID     string `json:"poi_id"`            // 关联的点位 ID
	ParentID  string `json:"parent_id"`         // 父评论 ID (空字符串表示一级评论)
	Content   string `json:"content"`           // 评论内容
	Status    int    `json:"status"`            // 0:待审, 1:通过, 2:拒绝
	LikeCount int    `json:"like_count"`        // 点赞数
	CreatedAt string `json:"created_at"`        // 业务创建时间
	UpdatedAt string `json:"updated_at"`        // 业务更新时间

	// 可选：如果前端展示需要，可在此处冗余 Nickname/Avatar，
	// 目前暂按 PRD 最小闭环设计，只存 OpenID。
}

// CommentQuery 评论筛选
type CommentQuery struct {
	POIID  string `form:"poi_id"`           // 查某个景点的评论
	Status int    `form:"status,default=1"` // 默认只查过审的
	Page   int    `form:"page,default=1"`
	Size   int    `form:"size,default=20"`
}
