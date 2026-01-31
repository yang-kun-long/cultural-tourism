// File: services/region_service.go
package services

import (
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

const CollectionRegion = "regions"

// CreateRegion 创建新区域
func CreateRegion(region models.Region) (map[string]interface{}, error) {
	// 补全默认值
	region.ID = "" // 安全置空，ID由云开发生成
	region.Status = 1
	region.CreatedAt = time.Now().Format(time.RFC3339)
	region.UpdatedAt = time.Now().Format(time.RFC3339)
	if region.Sort == 0 {
		region.Sort = 100 // 默认排序权重
	}

	result, err := tcb.Client.CreateData(CollectionRegion, region)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ListRegions 获取区域列表
func ListRegions(page, size, status int) (map[string]interface{}, error) {
	// 构造筛选条件
	where := map[string]interface{}{
		"status": map[string]interface{}{
			"$eq": status,
		},
	}

	filter := map[string]interface{}{
		"where": where,
		// TODO: 等待 SDK 支持 orderBy
	}

	result, err := tcb.Client.ListData(CollectionRegion, filter, page, size)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetRegionDetail 获取单条区域详情
func GetRegionDetail(id string) (map[string]interface{}, error) {
	result, err := tcb.Client.GetDetail(CollectionRegion, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateRegion 更新区域
func UpdateRegion(id string, region models.Region) error {
	// 使用 Map 构造更新数据，支持 Partial Update
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if region.Name != "" {
		updateData["name"] = region.Name
	}
	if region.Sort > 0 {
		updateData["sort"] = region.Sort
	}
	// 简单处理：只有非0才更新。如果需要更新为0，建议使用 map[string]interface{} 接收参数
	if region.Status != 0 {
		updateData["status"] = region.Status
	}

	err := tcb.Client.UpdateData(CollectionRegion, id, updateData)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRegion 删除区域
func DeleteRegion(id string) error {
	err := tcb.Client.DeleteData(CollectionRegion, id)
	if err != nil {
		return err
	}

	return nil
}