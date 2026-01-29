// File: controllers/region_controller.go
package controllers

import (
	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"
	"net/http"

	"github.com/gin-gonic/gin"
)

// =================================================================================
// 区域管理 (Regions) - 标准 REST API
// =================================================================================

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

	// 1. 绑定参数
	if err := c.ShouldBindJSON(&region); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	region.Status = 1 // 默认启用

	// 2. 调用云开发 HTTP API 写入数据
	result, err := tcb.Client.CreateData("regions", region)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "调用云开发失败",
			"detail": err.Error(),
		})
		return
	}

	// 3. 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    result,
	})
}

// GetRegions 获取区域列表
// @Summary      获取所有区域
// @Description  查询区域列表，支持分页（默认查前100条）
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /regions [get]
func GetRegions(c *gin.Context) {
	// 调用 ListData 方法
	// 【修正】增加 nil 参数，对应 ListData 新增的 filter 参数
	// 第 1 页，每页 100 条 (足够覆盖所有景区区域)
	result, err := tcb.Client.ListData("regions", nil, 1, 100)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "查询失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    result["data"], // 将腾讯云返回的 data 字段透传给前端
	})
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

	result, err := tcb.Client.GetDetail("regions", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "查询失败或记录不存在", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateRegion 更新区域
// @Summary      更新区域
// @Description  根据 ID 更新区域信息 (支持 name, sort, status)
// @Tags         Regions
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "区域ID"
// @Param        data  body      object  true  "更新内容(JSON)"
// @Success      200   {object}  map[string]interface{}
// @Router       /regions/{id} [put]
func UpdateRegion(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	// 定义一个只包含需要更新字段的结构体
	// 这里为了简单，我们允许更新 name, sort, status
	var updateData struct {
		Name   string `json:"name"`
		Sort   int    `json:"sort"`
		Status int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 调用更新
	err := tcb.Client.UpdateData("regions", id, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "id": id})
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
	// 1. 从 URL 路径参数获取 id
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	// 2. 调用云开发删除
	err := tcb.Client.DeleteData("regions", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "删除失败",
			"detail": err.Error(),
		})
		return
	}

	// 3. 返回成功
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
		"id":      id,
	})
}
