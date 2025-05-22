package controllers

import (
	"go-book/models"
	"go-book/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookController 处理图书相关的请求
type BookController struct {
	DB *gorm.DB
}

// NewBookController 创建新的图书控制器实例
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{DB: db}
}

// CreateBook 创建新图书
func (bc *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := bc.DB.Create(&book)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建图书失败"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBook 获取单本图书信息
func (bc *BookController) GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	result := bc.DB.First(&book, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// GetBooks 获取所有图书列表
func (bc *BookController) GetBooks(c *gin.Context) {
	var books []models.Book

	result := bc.DB.Find(&books)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图书列表失败"})
		return
	}

	c.JSON(http.StatusOK, books)
}

// UpdateBook 更新图书信息
func (bc *BookController) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var book models.Book
	if err := bc.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图书不存在"})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bc.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}

// DeleteBook 删除图书
func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	result := bc.DB.Delete(&models.Book{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除图书失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "图书已成功删除"})
}

// SearchBooks 搜索图书
func (bc *BookController) SearchBooks(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供搜索关键词"})
		return
	}

	// 创建豆瓣服务实例
	doubanService := services.NewDoubanService()

	// 搜索豆瓣图书
	books, err := doubanService.SearchBooks(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索图书失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}