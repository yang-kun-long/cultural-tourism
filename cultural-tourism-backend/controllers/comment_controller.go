// File: controllers/comment_controller.go
package controllers

import (
	"net/http"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/services"

	"github.com/gin-gonic/gin"
)


// CreateComment 发布评论
// @Summary      发布评论
// @Description  用户发布评论 (默认待审核 status=0)
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        comment  body      models.Comment  true  "评论信息"
// @Success      200      {object}  map[string]interface{}
// @Router       /comments [post]
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	result, err := services.CreateComment(&comment)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCommentList 获取评论列表
// @Summary      获取评论列表
// @Description  默认只显示审核通过(status=1)的评论
// @Tags         Comments
// @Param        poi_id  query  string  false  "点位ID"
// @Param        status  query  int     false  "状态 (1:通过, 0:待审)"
// @Param        page    query  int     false  "页码"
// @Param        size    query  int     false  "每页数量"
// @Success      200     {object}  map[string]interface{}
// @Router       /comments [get]
func GetCommentList(c *gin.Context) {
	var query models.CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListComments(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCommentDetail 获取评论详情
// @Summary      获取评论详情
// @Tags         Comments
// @Param        id   path      string  true  "评论ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /comments/{id} [get]
func GetCommentDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	result, err := services.GetCommentDetail(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateComment 更新评论 (审核/点赞)
// @Summary      更新评论 (审核/点赞)
// @Description  用于管理员审核 (修改status) 或 用户点赞 (修改like_count)
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "评论ID"
// @Param        comment  body      models.Comment true  "更新内容 (仅status/like_count)"
// @Success      200      {object}  map[string]interface{}
// @Router       /comments/{id} [put]
func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := services.UpdateComment(id, comment); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeleteComment 删除评论
// @Summary      删除评论
// @Description  管理员可删除违规评论，用户可删除自己的评论(依赖TCB权限)
// @Tags         Comments
// @Param        id   path      string  true  "评论ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteComment(id); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}
