package controllers

import (
	"cultural-tourism-backend/models"
	"cultural-tourism-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateFavorite 创建收藏
// @Summary 收藏资源
// @Description 用户收藏旅拍主题/景点POI/商品 (幂等：重复收藏返回错误)
// @Tags Favorites
// @Accept json
// @Produce json
// @Param body body models.FavoriteCreateRequest true "收藏请求"
// @Success 200 {object} map[string]interface{} "收藏成功"
// @Failure 400 {object} map[string]interface{} "参数错误或已收藏"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/favorites [post]
func CreateFavorite(c *gin.Context) {
	var req models.FavoriteCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	favorite := &models.Favorite{
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
	}

	result, err := services.CreateFavorite(favorite)
	if err != nil {
		// 已收藏错误返回 400
		if err.Error() == "already favorited" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already favorited", "code": "ALREADY_FAVORITED"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteFavorite 取消收藏
// @Summary 取消收藏
// @Description 通过资源类型和资源ID取消收藏 (RESTful路径参数)
// @Tags Favorites
// @Produce json
// @Param resource_type path string true "资源类型 (theme/poi/product)"
// @Param resource_id path string true "资源ID"
// @Success 200 {object} map[string]interface{} "取消成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 404 {object} map[string]interface{} "收藏记录不存在"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /favorites/{resource_type}/{resource_id} [delete]
func DeleteFavorite(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")

	if resourceType == "" || resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource_type and resource_id are required"})
		return
	}

	err := services.DeleteFavorite(resourceType, resourceID)
	if err != nil {
		if err.Error() == "favorite not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "favorite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfavorited successfully"})
}

// ListFavorites 获取收藏列表
// @Summary 获取用户收藏列表
// @Description 分页获取当前用户的收藏列表 (支持按资源类型筛选)
// @Tags Favorites
// @Produce json
// @Param resource_type query string false "资源类型筛选 (theme/poi/product)"
// @Param page query int false "页码 (默认1)" default(1)
// @Param size query int false "每页数量 (默认20, 最大100)" default(20)
// @Success 200 {object} map[string]interface{} "收藏列表"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /api/favorites [get]
func ListFavorites(c *gin.Context) {
	var req models.FavoriteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListFavorites(req.ResourceType, req.Page, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CheckFavoriteStatus 检查收藏状态
// @Summary 检查是否已收藏
// @Description 检查指定资源是否已被当前用户收藏 (RESTful路径参数)
// @Tags Favorites
// @Produce json
// @Param resource_type path string true "资源类型 (theme/poi/product)"
// @Param resource_id path string true "资源ID"
// @Success 200 {object} map[string]interface{} "返回 {is_favorited: true/false}"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "服务器错误"
// @Router /favorites/{resource_type}/{resource_id} [get]
func CheckFavoriteStatus(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceID := c.Param("resource_id")

	if resourceType == "" || resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource_type and resource_id are required"})
		return
	}

	isFavorited, err := services.CheckFavoriteStatus(resourceType, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorited": isFavorited})
}
