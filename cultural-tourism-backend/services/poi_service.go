// File: services/poi_service.go
package services

import (
	"fmt"
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
)

// Collection name for POIs
const CollectionPOI = "poi"

// CreatePOI creates a new POI
func CreatePOI(poi *models.POI) (map[string]interface{}, error) {
	// [Security] 强制初始化字段，防止恶意篡改
	poi.ID = ""
	poi.Status = 1
	poi.CreatedAt = time.Now().Format(time.RFC3339)
	poi.UpdatedAt = time.Now().Format(time.RFC3339)

	// Images 数组防止 nil 导致入库异常
	if len(poi.Images) == 0 {
		poi.Images = []string{}
	}

	// 防止恶意写入非预期字段（业务兜底）


	return tcb.Client.CreateData(CollectionPOI, poi)
}

// ListPOIs retrieves POI list with filtering and pagination
func ListPOIs(query models.POIQuery) (map[string]interface{}, error) {
	where := make(map[string]interface{})

	// 状态筛选 - 默认只返回上线状态
	where["status"] = map[string]interface{}{"$eq": 1}

	// 区域筛选
	if query.RegionID != "" {
		where["region_id"] = map[string]interface{}{"$eq": query.RegionID}
	}

	// POI 类型筛选
	if query.Type != "" {
		where["type"] = map[string]interface{}{"$eq": query.Type}
	}



	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionPOI, filter, query.Page, query.Size)
}

// GetPOIDetail retrieves a single POI by ID
func GetPOIDetail(id string) (map[string]interface{}, error) {
	return tcb.Client.GetDetail(CollectionPOI, id)
}

// UpdatePOI updates an existing POI
func UpdatePOI(id string, poi *models.POI) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}

	if poi.Name != "" {
		updateData["name"] = poi.Name
	}
	if poi.Type != "" {
		updateData["type"] = poi.Type
	}
	if poi.RegionID != "" {
		updateData["region_id"] = poi.RegionID
	}
	if poi.Latitude != 0 {
		updateData["latitude"] = poi.Latitude
	}
	if poi.Longitude != 0 {
		updateData["longitude"] = poi.Longitude
	}
	if len(poi.Images) > 0 {
		updateData["images"] = poi.Images
	}
	if poi.Desc != "" {
		updateData["desc"] = poi.Desc
	}
	if poi.Address != "" {
		updateData["address"] = poi.Address
	}
	if poi.Phone != "" {
		updateData["phone"] = poi.Phone
	}
	if poi.OpenTime != "" {
		updateData["open_time"] = poi.OpenTime
	}

	if poi.Status != 0 {
		updateData["status"] = poi.Status
	}

	return tcb.Client.UpdateData(CollectionPOI, id, updateData)
}

// DeletePOI deletes a POI by ID
func DeletePOI(id string) error {
	return tcb.Client.DeleteData(CollectionPOI, id)
}

// GetNearbyPOIs retrieves POIs within a radius of a point
func GetNearbyPOIs(lat, lng float64, radiusKM int, page, size int) (map[string]interface{}, error) {
	where := map[string]interface{}{
		"status": map[string]interface{}{"$eq": 1}, // 只返回上线 POI
		"location": map[string]interface{}{
			"$near": map[string]interface{}{
				"$geometry": map[string]interface{}{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"$maxDistance": radiusKM * 1000,
			},
		},
	}

	filter := map[string]interface{}{
		"where": where,
	}

	return tcb.Client.ListData(CollectionPOI, filter, page, size)
}

// CountPOIsByRegion counts POIs by region (for statistics)
func CountPOIsByRegion(regionID string) (map[string]interface{}, error) {
	where := map[string]interface{}{
		"region_id": map[string]interface{}{"$eq": regionID},
		"status":    map[string]interface{}{"$eq": 1},
	}

	filter := map[string]interface{}{
		"where": where,
		"count": true,
	}

	return tcb.Client.ListData(CollectionPOI, filter, 1, 1)
}

// BatchUpdatePOIStatus batch updates POI status (for admin operations)
func BatchUpdatePOIStatus(ids []string, status int) error {
	for _, id := range ids {
		updateData := map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().Format(time.RFC3339),
		}
		if err := tcb.Client.UpdateData(CollectionPOI, id, updateData); err != nil {
			return fmt.Errorf("failed to update POI %s: %v", id, err)
		}
	}
	return nil
}