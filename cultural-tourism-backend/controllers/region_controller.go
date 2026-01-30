// File: controllers/region_controller.go
package controllers

import (
	"net/http"
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"

	"github.com/gin-gonic/gin"
)

const CollectionRegion = "regions"

// CreateRegion 创建新区域
// @Summary      创建区域
// @Description  创建一个新的景区区域 (Name必填)
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        region  body      models.Region  true  "区域信息"
// @Success      200     {object}  map[string]interface{}
// @Router       /regions [post]
func CreateRegion(c *gin.Context) {
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRegions 获取区域列表
// @Summary      获取所有区域
// @Description  查询区域列表，支持分页
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        page    query     int     false  "页码 (默认1)"
// @Param        size    query     int     false  "每页数量 (默认100)"
// @Param        status  query     int     false  "状态 (1:启用, 0:禁用)"
// @Success      200     {object}  map[string]interface{}
// @Router       /regions [get]
func GetRegions(c *gin.Context) {
	// 获取分页参数
	type Query struct {
		Page   int `form:"page,default=1"`
		Size   int `form:"size,default=100"`
		Status int `form:"status,default=1"`
	}
	var query Query
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 构造筛选条件
	where := map[string]interface{}{
		"status": map[string]interface{}{
			"$eq": query.Status,
		},
	}

	filter := map[string]interface{}{
		"where": where,
		// TODO: 等待 SDK 支持 orderBy
	}

	// 2. 调用 SDK
	result, err := tcb.Client.ListData(CollectionRegion, filter, query.Page, query.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRegionDetail 获取单条区域详情
// @Summary      获取区域详情
// @Description  根据 ID 获取单个区域的详细信息
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "区域ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /regions/{id} [get]
func GetRegionDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	result, err := tcb.Client.GetDetail(CollectionRegion, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "区域不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateRegion 更新区域
// @Summary      更新区域
// @Description  根据 ID 更新区域信息 (支持 name, sort, status)
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        id    path      string         true  "区域ID"
// @Param        data  body      models.Region  true  "更新内容 (仅需传修改字段)"
// @Success      200   {object}  map[string]interface{}
// @Router       /regions/{id} [put]
func UpdateRegion(c *gin.Context) {
	id := c.Param("id")
	var region models.Region
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// [Security Fix] 使用 Map 构造更新数据，支持 Partial Update
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeleteRegion 删除区域
// @Summary      删除区域
// @Description  根据 ID 删除指定区域
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "区域ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /regions/{id} [delete]
func DeleteRegion(c *gin.Context) {
	id := c.Param("id")
	err := tcb.Client.DeleteData(CollectionRegion, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
