// File: services/theme_service.go
package services

import (
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

// Collection name for Themes
const CollectionTheme = "theme"

// CreateTheme creates a new theme
func CreateTheme(theme *models.Theme) (map[string]interface{}, error) {
	// [Security] 强制初始化字段，防止恶意篡改
	theme.ID = ""
	theme.Status = 1
	theme.CreatedAt = time.Now().Format(time.RFC3339)
	theme.UpdatedAt = time.Now().Format(time.RFC3339)

	// 防止恶意写入非预期字段（业务兜底）
	if theme.Sort <= 0 {
		theme.Sort = 9999
	}

	return tcb.Client.CreateData(CollectionTheme, theme)
}

// ListThemes retrieves theme list with filtering and pagination
func ListThemes(query models.ThemeQuery) (map[string]interface{}, error) {
	where := make(map[string]interface{})

	// 状态筛选 - 默认只返回上线状态
	where["status"] = map[string]interface{}{"$eq": 1}

	// 区域筛选（区域优先）
	if query.RegionID != "" {
		where["region_id"] = map[string]interface{}{"$eq": query.RegionID}
	}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionTheme, filter, query.Page, query.Size)
}

// GetThemeDetail retrieves a single theme by ID
func GetThemeDetail(id string) (map[string]interface{}, error) {
	return tcb.Client.GetDetail(CollectionTheme, id)
}

// UpdateTheme updates an existing theme
func UpdateTheme(id string, theme *models.Theme) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if theme.Name != "" {
		updateData["name"] = theme.Name
	}
	if theme.Cover != "" {
		updateData["cover"] = theme.Cover
	}
	if theme.Desc != "" {
		updateData["desc"] = theme.Desc
	}
	if theme.RegionID != "" {
		updateData["region_id"] = theme.RegionID
	}
	if theme.Sort != 0 {
		updateData["sort"] = theme.Sort
	}
	if theme.Status != 0 {
		updateData["status"] = theme.Status
	}

	return tcb.Client.UpdateData(CollectionTheme, id, updateData)
}

// DeleteTheme deletes a theme by ID
func DeleteTheme(id string) error {
	return tcb.Client.DeleteData(CollectionTheme, id)
}

// GetThemesByRegion retrieves all themes for a specific region
func GetThemesByRegion(regionID string, page, size int) (map[string]interface{}, error) {
	where := map[string]interface{}{
		"region_id": map[string]interface{}{"$eq": regionID},
		"status":    map[string]interface{}{"$eq": 1},
	}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionTheme, filter, page, size)
}

// BatchUpdateThemeStatus batch updates theme status (for admin operations)
func BatchUpdateThemeStatus(ids []string, status int) error {
	for _, id := range ids {
		updateData := map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().Format(time.RFC3339),
		}
		if err := tcb.Client.UpdateData(CollectionTheme, id, updateData); err != nil {
			return err
		}
	}
	return nil
}

// CountThemesByRegion counts themes by region (for statistics)
func CountThemesByRegion(regionID string) (map[string]interface{}, error) {
	where := map[string]interface{}{
		"region_id": map[string]interface{}{"$eq": regionID},
		"status":    map[string]interface{}{"$eq": 1},
	}

	filter := map[string]interface{}{
		"where": where,
		"count": true,
	}

	return tcb.Client.ListData(CollectionTheme, filter, 1, 1)
}