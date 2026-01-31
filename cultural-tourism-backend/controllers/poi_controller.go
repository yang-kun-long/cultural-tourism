// File: controllers/poi_controller.go
package controllers

import (
	"math"
	"net/http"
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"

	"github.com/gin-gonic/gin"
)

const CollectionPOI = "pois"

// ... (calculateDistance 函数保持不变，此处省略以节省篇幅，请保留之前的实现) ...
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // 地球半径 (米)
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// CreatePOI 创建点位
// @Summary      创建点位
// @Description  创建新的 POI 点位 (景点/饭店/酒店/旅拍机)
// @Tags         POI
// @Accept       json
// @Produce      json
// @Param        poi  body      models.POI  true  "POI信息"
// @Success      200  {object}  map[string]interface{}
// @Router       /pois [post]
func CreatePOI(c *gin.Context) {
	var poi models.POI
	if err := c.ShouldBindJSON(&poi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// [Security] 强制覆盖系统字段
	poi.ID = ""
	poi.CreatedAt = time.Now().Format(time.RFC3339)
	poi.UpdatedAt = time.Now().Format(time.RFC3339)
	poi.Distance = 0
	if poi.Status == 0 {
		poi.Status = 1
	}

	result, err := tcb.Client.CreateData(CollectionPOI, poi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create POI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPOIList 获取点位列表
// @Summary      获取点位列表 (支持LBS)
// @Description  查询点位列表，支持按区域、类型筛选。若传入 lat/lng，结果将包含距离信息(_distance)。
// @Tags         POI
// @Param        region_id  query  string   false  "区域ID"
// @Param        type       query  string   false  "点位类型 (scenic/food/hotel/booth)"
// @Param        lat        query  float64  false  "用户纬度 (用于计算距离)"
// @Param        lng        query  float64  false  "用户经度 (用于计算距离)"
// @Param        page       query  int      false  "页码"
// @Param        size       query  int      false  "每页数量"
// @Success      200        {object} map[string]interface{}
// @Router       /pois [get]
func GetPOIList(c *gin.Context) {
	var query models.POIQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 构造 TCB 标准 Filter
	where := make(map[string]interface{})
	where["status"] = map[string]interface{}{"$eq": 1}

	if query.RegionID != "" {
		where["region_id"] = map[string]interface{}{"$eq": query.RegionID}
	}
	if query.Type != "" {
		where["type"] = map[string]interface{}{"$eq": query.Type}
	}

	filter := map[string]interface{}{
		"where": where,
	}

	// 2. 调用 SDK
	result, err := tcb.Client.ListData(CollectionPOI, filter, query.Page, query.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch POIs: " + err.Error()})
		return
	}

	// 3. [LBS Feature] 距离计算
	if query.UserLat != 0 && query.UserLng != 0 {
		if dataMap, ok := result["data"].(map[string]interface{}); ok {
			if records, ok := dataMap["records"].([]interface{}); ok {
				for i, r := range records {
					if rec, ok := r.(map[string]interface{}); ok {
						lat, _ := rec["latitude"].(float64)
						lng, _ := rec["longitude"].(float64)

						if lat != 0 && lng != 0 {
							dist := calculateDistance(query.UserLat, query.UserLng, lat, lng)
							rec["_distance"] = math.Round(dist)
							records[i] = rec
						}
					}
				}
				dataMap["records"] = records
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetPOI 获取单个点位详情
// @Summary      获取点位详情
// @Tags         POI
// @Param        id   path      string  true  "POI ID"
// @Success      200  {object}  models.POI
// @Router       /pois/{id} [get]
func GetPOI(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	result, err := services.GetPOIDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POI not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePOI 更新点位
// @Summary      更新点位
// @Description  更新点位信息 (安全更新，忽略系统字段)
// @Tags         POI
// @Accept       json
// @Produce      json
// @Param        id   path      string      true  "POI ID"
// @Param        poi  body      models.POI  true  "更新信息"
// @Success      200  {object}  map[string]interface{}
// @Router       /pois/{id} [put]
func UpdatePOI(c *gin.Context) {
	id := c.Param("id")
	var poi models.POI
	if err := c.ShouldBindJSON(&poi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// [Security Fix] Partial Update Map
	updateData := map[string]interface{}{
		"updated_at": time.Now().Format(time.RFC3339),
	}
	// ... (白名单字段赋值逻辑同前文，此处省略以节省篇幅，请保留) ...
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

	err := tcb.Client.UpdateData(CollectionPOI, id, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update POI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeletePOI 删除点位
// @Summary      删除点位
// @Tags         POI
// @Param        id   path      string  true  "POI ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /pois/{id} [delete]
func DeletePOI(c *gin.Context) {
	id := c.Param("id")
	err := tcb.Client.DeleteData(CollectionPOI, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete POI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
