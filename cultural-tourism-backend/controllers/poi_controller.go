// File: controllers/poi_controller.go
package controllers

import (
	"net/http"
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"

	"github.com/gin-gonic/gin"
)

const CollectionPOI = "pois"

// CreatePOI 创建点位
// @Summary 创建点位
// @Tags POI
// @Accept json
// @Produce json
// @Param poi body models.POI true "POI Info"
// @Success 200 {object} map[string]interface{}
// @Router /pois [post]
func CreatePOI(c *gin.Context) {
	var poi models.POI
	if err := c.ShouldBindJSON(&poi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 补全基础字段
	poi.CreatedAt = time.Now().Format(time.RFC3339)
	poi.UpdatedAt = time.Now().Format(time.RFC3339)
	if poi.Status == 0 {
		poi.Status = 1 // 默认启用
	}

	result, err := tcb.Client.CreateData(CollectionPOI, poi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create POI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPOIList 获取点位列表
// @Summary 获取点位列表
// @Tags POI
// @Param region_id query string false "区域ID"
// @Param type query string false "点位类型"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} map[string]interface{}
// @Router /pois [get]
func GetPOIList(c *gin.Context) {
	var query models.POIQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建筛选条件
	filter := make(map[string]interface{})
	if query.RegionID != "" {
		filter["region_id"] = query.RegionID
	}
	if query.Type != "" {
		filter["type"] = query.Type
	}

	// 调用 Client.ListData (已包含 filter 参数)
	result, err := tcb.Client.ListData(CollectionPOI, filter, query.Page, query.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch POIs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPOI 获取单个点位详情
// @Summary 获取点位详情
// @Tags POI
// @Param id path string true "POI ID"
// @Success 200 {object} models.POI
// @Router /pois/{id} [get]
func GetPOI(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	result, err := tcb.Client.GetDetail(CollectionPOI, id)
	if err != nil {
		if err.Error() == "未找到该记录" {
			c.JSON(http.StatusNotFound, gin.H{"error": "POI not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get POI: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePOI 更新点位
// @Summary 更新点位
// @Tags POI
// @Accept json
// @Produce json
// @Param id path string true "POI ID"
// @Param poi body models.POI true "Update Info"
// @Success 200 {object} map[string]interface{}
// @Router /pois/{id} [put]
func UpdatePOI(c *gin.Context) {
	id := c.Param("id")
	var poi models.POI
	if err := c.ShouldBindJSON(&poi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	poi.UpdatedAt = time.Now().Format(time.RFC3339)

	err := tcb.Client.UpdateData(CollectionPOI, id, poi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update POI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeletePOI 删除点位
// @Summary 删除点位
// @Tags POI
// @Param id path string true "POI ID"
// @Success 200 {object} map[string]interface{}
// @Router /pois/{id} [delete]
func DeletePOI(c *gin.Context) {
	id := c.Param("id")

	err := tcb.Client.DeleteData(CollectionPOI, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete POI: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
