// File: controllers/photo_controller.go
package controllers

import (
	"net/http"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/services"

	"github.com/gin-gonic/gin"
)


// CreatePhoto 上传照片
// @Summary      上传照片
// @Description  用户上传旅拍照片 (上传后默认为待审核状态 status=0)
// @Tags         Photos
// @Accept       json
// @Produce      json
// @Param        photo  body      models.Photo  true  "照片信息"
// @Success      200    {object}  map[string]interface{}
// @Router       /photos [post]
func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	result, err := services.CreatePhoto(&photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPhotoList 获取照片列表 (瀑布流)
// @Summary      获取照片列表
// @Description  支持按主题ID筛选，默认只显示审核通过(status=1)的照片
// @Tags         Photos
// @Param        theme_id  query  string  false  "主题ID"
// @Param        status    query  int     false  "状态 (1:通过, 0:待审)"
// @Param        page      query  int     false  "页码"
// @Param        size      query  int     false  "每页数量"
// @Success      200       {object}  map[string]interface{}
// @Router       /photos [get]
func GetPhotoList(c *gin.Context) {
	var query models.PhotoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListPhotos(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPhotoDetail 获取照片详情
// @Summary      获取照片详情
// @Tags         Photos
// @Param        id   path      string  true  "照片ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /photos/{id} [get]
func GetPhotoDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	result, err := services.GetPhotoDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "照片不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePhoto 更新照片状态 (审核/点赞)
// @Summary      更新照片 (审核/点赞)
// @Description  用于管理员审核 (修改status) 或 用户点赞 (修改like_count)
// @Tags         Photos
// @Accept       json
// @Produce      json
// @Param        id     path      string        true  "照片ID"
// @Param        photo  body      models.Photo  true  "更新内容 (仅status/like_count)"
// @Success      200    {object}  map[string]interface{}
// @Router       /photos/{id} [put]
func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := services.UpdatePhoto(id, photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeletePhoto 删除照片
// @Summary      删除照片
// @Description  管理员可删除违规照片，用户可删除自己的照片(依赖TCB权限)
// @Tags         Photos
// @Param        id   path      string  true  "照片ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /photos/{id} [delete]
func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	err := services.DeletePhoto(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
