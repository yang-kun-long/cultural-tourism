// File: routes/router.go
package routes

import (
	"cultural-tourism-backend/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	// CORS 配置
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// ==============================
		// Regions 区域管理 (标准 REST API)
		// ==============================

		// 1. 创建 (Create)
		api.POST("/regions", controllers.CreateRegion)

		// 2. 列表查询 (Read List)
		api.GET("/regions", controllers.GetRegions)

		// 3. 单条详情 (Read Detail) - 新增
		api.GET("/regions/:id", controllers.GetRegionDetail)

		// 4. 更新 (Update) - 新增
		// 使用 PUT 用于全量/部分更新
		api.PUT("/regions/:id", controllers.UpdateRegion)

		// 5. 删除 (Delete)
		api.DELETE("/regions/:id", controllers.DeleteRegion)

		// === Phase 3: POI (点位管理) ===
		api.POST("/pois", controllers.CreatePOI)       // 创建
		api.GET("/pois", controllers.GetPOIList)       // 列表 (支持 region_id, type 筛选)
		api.GET("/pois/:id", controllers.GetPOI)       // 详情
		api.PUT("/pois/:id", controllers.UpdatePOI)    // 更新
		api.DELETE("/pois/:id", controllers.DeletePOI) // 删除

		// ================= Phase 4: UGC 旅拍主题 (Themes) =================
		api.POST("/themes", controllers.CreateTheme)       // 创建主题
		api.GET("/themes", controllers.GetThemeList)       // 列表 (支持 ?region_id=...)
		api.GET("/themes/:id", controllers.GetThemeDetail) // 详情
		api.PUT("/themes/:id", controllers.UpdateTheme)    // 更新
		api.DELETE("/themes/:id", controllers.DeleteTheme) // 删除

		// ================= Phase 4 (Part 2): UGC 照片管理 (Photos) =================
		api.POST("/photos", controllers.CreatePhoto)       // 上传 (默认待审)
		api.GET("/photos", controllers.GetPhotoList)       // 瀑布流 (默认查已过审)
		api.GET("/photos/:id", controllers.GetPhotoDetail) // 详情
		api.PUT("/photos/:id", controllers.UpdatePhoto)    // 审核/点赞
		api.DELETE("/photos/:id", controllers.DeletePhoto) // 删除

		// ================= Phase 5: 评论互动 (Comments) =================
		api.POST("/comments", controllers.CreateComment)       // 发布评论
		api.GET("/comments", controllers.GetCommentList)       // 列表
		api.GET("/comments/:id", controllers.GetCommentDetail) // 详情
		api.PUT("/comments/:id", controllers.UpdateComment)    // 审核/点赞
		api.DELETE("/comments/:id", controllers.DeleteComment) // 删除

		// ================= Phase 5: 商品导流 (Products) =================
		api.POST("/products", controllers.CreateProduct)       // 创建
		api.GET("/products", controllers.GetProductList)       // 列表
		api.GET("/products/:id", controllers.GetProductDetail) // 详情
		api.PUT("/products/:id", controllers.UpdateProduct)    // 更新
		api.DELETE("/products/:id", controllers.DeleteProduct) // 删除

		// ================= Phase 6: 收藏体系 (Favorites) =================
		api.POST("/favorites", controllers.CreateFavorite)                                    // 收藏资源
		api.DELETE("/favorites/:resource_type/:resource_id", controllers.DeleteFavorite)      // 取消收藏 (RESTful)
		api.GET("/favorites", controllers.ListFavorites)                                      // 收藏列表
		api.GET("/favorites/:resource_type/:resource_id", controllers.CheckFavoriteStatus)    // 检查收藏状态
	}
}


