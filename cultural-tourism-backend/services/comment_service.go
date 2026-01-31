// File: services/comment_service.go
package services

import (
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

// [Critical] 数据库实际集合名为单数 "comment"
const CollectionComment = "comment"

// CreateComment 创建评论（默认待审）
func CreateComment(comment *models.Comment) (map[string]interface{}, error) {
	comment.ID = ""
	comment.Status = 0
	comment.LikeCount = 0
	comment.CreatedAt = time.Now().Format(time.RFC3339)
	comment.UpdatedAt = time.Now().Format(time.RFC3339)
	if comment.ParentID == "" {
		comment.ParentID = ""
	}

	return tcb.Client.CreateData(CollectionComment, comment)
}

// ListComments 获取评论列表
func ListComments(query models.CommentQuery) (map[string]interface{}, error) {
	where := make(map[string]interface{})
	where["status"] = map[string]interface{}{"$eq": query.Status}
	if query.POIID != "" {
		where["poi_id"] = map[string]interface{}{"$eq": query.POIID}
	}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionComment, filter, query.Page, query.Size)
}

// GetCommentDetail 获取评论详情
func GetCommentDetail(id string) (map[string]interface{}, error) {
	return tcb.Client.GetDetail(CollectionComment, id)
}

// UpdateComment 更新评论（审核/点赞）
func UpdateComment(id string, comment models.Comment) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if comment.Status != 0 {
		updateData["status"] = comment.Status
	}
	if comment.LikeCount > 0 {
		updateData["like_count"] = comment.LikeCount
	}

	return tcb.Client.UpdateData(CollectionComment, id, updateData)
}

// DeleteComment 删除评论
func DeleteComment(id string) error {
	return tcb.Client.DeleteData(CollectionComment, id)
}
