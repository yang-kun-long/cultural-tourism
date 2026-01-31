// File: services/photo_service.go
package services

import (
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

// [Critical] 数据库实际集合名为单数 "photo"
const CollectionPhoto = "photo"

// CreatePhoto 上传照片（默认待审）
func CreatePhoto(photo *models.Photo) (map[string]interface{}, error) {
	photo.ID = ""
	photo.Status = 0
	photo.LikeCount = 0
	photo.CreatedAt = time.Now().Format(time.RFC3339)
	photo.UpdatedAt = time.Now().Format(time.RFC3339)

	return tcb.Client.CreateData(CollectionPhoto, photo)
}

// ListPhotos 获取照片列表
func ListPhotos(query models.PhotoQuery) (map[string]interface{}, error) {
	where := make(map[string]interface{})
	where["status"] = map[string]interface{}{"$eq": query.Status}
	if query.ThemeID != "" {
		where["theme_id"] = map[string]interface{}{"$eq": query.ThemeID}
	}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionPhoto, filter, query.Page, query.Size)
}

// GetPhotoDetail 获取照片详情
func GetPhotoDetail(id string) (map[string]interface{}, error) {
	return tcb.Client.GetDetail(CollectionPhoto, id)
}

// UpdatePhoto 更新照片（审核/点赞）
func UpdatePhoto(id string, photo models.Photo) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if photo.Status != 0 {
		updateData["status"] = photo.Status
	}
	if photo.LikeCount > 0 {
		updateData["like_count"] = photo.LikeCount
	}

	return tcb.Client.UpdateData(CollectionPhoto, id, updateData)
}

// DeletePhoto 删除照片
func DeletePhoto(id string) error {
	return tcb.Client.DeleteData(CollectionPhoto, id)
}
