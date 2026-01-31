// File: controllers/product_controller.go
package controllers

import (
	"net/http"

	"cultural-tourism-backend/models"
	"cultural-tourism-backend/services"

	"github.com/gin-gonic/gin"
)

// CreateProduct 创建商品导流
// @Summary      创建商品导流
// @Description  创建商品导流信息 (无支付，仅跳转)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "商品信息"
// @Success      200      {object}  map[string]interface{}
// @Router       /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	result, err := services.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProductList 获取商品列表
// @Summary      获取商品列表
// @Description  商品导流列表（无支付，点击跳转）
// @Tags         Products
// @Param        page  query  int  false  "页码"
// @Param        size  query  int  false  "每页数量"
// @Success      200   {object}  map[string]interface{}
// @Router       /products [get]
func GetProductList(c *gin.Context) {
	var query models.ProductQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.ListProducts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProductDetail 获取商品详情
// @Summary      获取商品详情
// @Tags         Products
// @Param        id   path      string  true  "商品ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func GetProductDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID不能为空"})
		return
	}

	result, err := services.GetProductDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateProduct 更新商品导流
// @Summary      更新商品导流
// @Description  更新商品导流信息 (支持增量更新)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "商品ID"
// @Param        product  body      models.Product true  "更新内容"
// @Success      200      {object}  map[string]interface{}
// @Router       /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := services.UpdateProduct(id, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}

// DeleteProduct 删除商品导流
// @Summary      删除商品导流
// @Tags         Products
// @Param        id   path      string  true  "商品ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "id": id})
}