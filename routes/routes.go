package routes

import (
	"go-book/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建图书控制器实例
	bookController := controllers.NewBookController(db)

	// 图书相关路由
	api := r.Group("/api")

	api.GET("/ping", controllers.NewCommonController().Ping)

	{
		// 图书管理路由
		books := api.Group("/books")
		{
			books.POST("/", bookController.CreateBook)      // 创建图书
			books.GET("/", bookController.GetBooks)         // 获取所有图书
			books.GET("/sq/search", bookController.SearchBooks)  // 搜索图书
			books.GET("/:id", bookController.GetBook)       // 获取单本图书
			books.PUT("/:id", bookController.UpdateBook)    // 更新图书
			books.DELETE("/:id", bookController.DeleteBook) // 删除图书
		}
	}
}
