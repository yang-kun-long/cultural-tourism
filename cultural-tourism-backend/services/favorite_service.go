package services

import (
	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
	"errors"
	"time"
)

const CollectionFavorites = "favorites"

// CreateFavorite 创建收藏 (幂等：重复收藏同一资源会返回已存在错误)
func CreateFavorite(favorite *models.Favorite) (map[string]interface{}, error) {
	// 安全处理：剥离系统字段
	favorite.ID = ""
	favorite.OpenID = "" // 由 TCB 自动注入

	// 设置时间戳
	favorite.CreatedAt = time.Now().Format(time.RFC3339)
	favorite.UpdatedAt = time.Now().Format(time.RFC3339)

	// 验证资源类型
	validTypes := map[string]bool{"theme": true, "poi": true, "product": true}
	if !validTypes[favorite.ResourceType] {
		return nil, errors.New("invalid resource_type, must be one of: theme, poi, product")
	}

	// 检查是否已收藏 (防止重复)
	existFilter := map[string]interface{}{
		"where": map[string]interface{}{
			"resource_type": map[string]interface{}{"$eq": favorite.ResourceType},
			"resource_id":   map[string]interface{}{"$eq": favorite.ResourceID},
		},
	}

	existResult, err := tcb.Client.ListData(CollectionFavorites, existFilter, 1, 1)
	if err != nil {
		return nil, err
	}

	// 检查返回的数据结构
	if dataMap, ok := existResult["data"].(map[string]interface{}); ok {
		if records, ok := dataMap["records"].([]interface{}); ok && len(records) > 0 {
			return nil, errors.New("already favorited")
		}
	}

	// 创建收藏记录
	return tcb.Client.CreateData(CollectionFavorites, favorite)
}

// DeleteFavorite 取消收藏 (通过资源类型和资源ID删除)
func DeleteFavorite(resourceType, resourceID string) error {
	// 验证资源类型
	validTypes := map[string]bool{"theme": true, "poi": true, "product": true}
	if !validTypes[resourceType] {
		return errors.New("invalid resource_type, must be one of: theme, poi, product")
	}

	// 构造查询条件 - 先找到要删除的记录
	filter := map[string]interface{}{
		"where": map[string]interface{}{
			"resource_type": map[string]interface{}{"$eq": resourceType},
			"resource_id":   map[string]interface{}{"$eq": resourceID},
		},
	}

	// 先查询记录获取 _id
	result, err := tcb.Client.ListData(CollectionFavorites, filter, 1, 1)
	if err != nil {
		return err
	}

	// 检查是否找到记录
	if dataMap, ok := result["data"].(map[string]interface{}); ok {
		if records, ok := dataMap["records"].([]interface{}); ok && len(records) > 0 {
			if record, ok := records[0].(map[string]interface{}); ok {
				if id, ok := record["_id"].(string); ok {
					// 使用 _id 删除
					return tcb.Client.DeleteData(CollectionFavorites, id)
				}
			}
		}
	}

	return errors.New("favorite not found")
}

// ListFavorites 获取用户收藏列表 (支持资源类型筛选和分页)
func ListFavorites(resourceType string, page, size int) (map[string]interface{}, error) {
	// 设置默认分页
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	// 构造查询条件
	where := map[string]interface{}{}
	if resourceType != "" {
		// 验证资源类型
		validTypes := map[string]bool{"theme": true, "poi": true, "product": true}
		if !validTypes[resourceType] {
			return nil, errors.New("invalid resource_type, must be one of: theme, poi, product")
		}
		where["resource_type"] = map[string]interface{}{"$eq": resourceType}
	}

	filter := map[string]interface{}{
		"where":   where,
		"orderBy": []interface{}{
			map[string]interface{}{"field": "created_at", "order": "desc"},
		},
	}

	return tcb.Client.ListData(CollectionFavorites, filter, page, size)
}

// CheckFavoriteStatus 检查收藏状态 (用于前端判断是否已收藏)
func CheckFavoriteStatus(resourceType, resourceID string) (bool, error) {
	// 验证资源类型
	validTypes := map[string]bool{"theme": true, "poi": true, "product": true}
	if !validTypes[resourceType] {
		return false, errors.New("invalid resource_type, must be one of: theme, poi, product")
	}

	filter := map[string]interface{}{
		"where": map[string]interface{}{
			"resource_type": map[string]interface{}{"$eq": resourceType},
			"resource_id":   map[string]interface{}{"$eq": resourceID},
		},
	}

	result, err := tcb.Client.ListData(CollectionFavorites, filter, 1, 1)
	if err != nil {
		return false, err
	}

	// 检查返回的数据结构
	if dataMap, ok := result["data"].(map[string]interface{}); ok {
		if records, ok := dataMap["records"].([]interface{}); ok && len(records) > 0 {
			return true, nil
		}
	}

	return false, nil
}
