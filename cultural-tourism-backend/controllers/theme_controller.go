// File: controllers/theme_controller.go
package controllers

import (
	"net/http"
	"time"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/tcb"

	"github.com/gin-gonic/gin"
)

const CollectionTheme = "theme" // 对应 model-json 中的 name

// CreateTheme 创建旅拍主题
// @Summary      创建旅拍主题
// @Description  创建新的旅拍活动主题 (仅管理员)
// @Tags         Themes
// @Accept       json
// @Produce      json
// @Param        theme  body      models.Theme  true  "主题信息"
// @Success      200    {object}  map[string]interface{}
// @Router       /themes [post]
func CreateTheme(c *gin.Context) {
	var theme models.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 补全默认值
	theme.ID = ""
	theme.CreatedAt = time.Now().Format(time.RFC3339)
	theme.UpdatedAt = time.Now().Format(time.RFC3339)
	if theme.Sort == 0 {
		theme.Sort = 100
	}
	if theme.Status == 0 {
		theme.Status = 1 // 默认上架
	}

	result, err := tcb.Client.CreateData(CollectionTheme, theme)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetThemeList 获取主题列表
// @Summary      获取主题列表
// @Description  支持按区域筛选。实现PRD“区域优先推荐”：前端应先传region_id查询，若为空则不传region_id查全局。
// @Tags         Themes
// @Param        region_id  query  string  false  "区域ID"
// @Param        status     query  int     false  "状态 (1:启用)"
// @Param        page       query  int     false  "页码"
// @Param        size       query  int     false  "每页数量"
// @Success      200        {object}  map[string]interface{}
// @Router       /themes [get]
func GetThemeList(c *gin.Context) {
	var query models.ThemeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 构造筛选条件
	where := map[string]interface{}{
		"status": map[string]interface{}{"$eq": query.Status},
	}

	if query.RegionID != "" {
		where["region_id"] = map[string]interface{}{"$eq": query.RegionID}
	}

	filter := map[string]interface{}{
		"where": where,
		// TODO: 待 SDK 支持 orderBy 后，此处应增加 "orderBy": [{"field": "sort", "direction": "desc"}]
	}

	// 2. 调用 SDK
	result, err := tcb.Client.ListData(CollectionTheme, filter, query.Page, query.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetThemeDetail 获取主题详情
// @Summary      获取主题详情
// @Tags         Themes
// @Param        id   path      string  true  "主题ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /themes/{id} [get]
func GetThemeDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	result, err := tcb.Client.GetDetail(CollectionTheme, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主题不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateTheme 更新主题
// @Summary      更新主题
// @Description  更新主题信息 (支持增量更新)
// @Tags         Themes
// @Accept       json
// @Produce      json
// @Param        id     path      string        true  "主题ID"
// @Param        theme  body      models.Theme  true  "更新内容"
// @Success      200    {object}  map[string]interface{}
// @Router       /themes/{id} [put]
func UpdateTheme(c *gin.Context) {
	id := c.Param("id")
	var theme models.Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// [Security] Partial Update Map 构造
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
	if theme.Sort > 0 {
		updateData["sort"] = theme.Sort
	}
	if theme.Status != 0 {
		updateData["status"] = theme.Status
	}

	err := tcb.Client.UpdateData(CollectionTheme, id, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeleteTheme 删除主题
// @Summary      删除主题
// @Tags         Themes
// @Param        id   path      string  true  "主题ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /themes/{id} [delete]
func DeleteTheme(c *gin.Context) {
	id := c.Param("id")
	err := tcb.Client.DeleteData(CollectionTheme, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
