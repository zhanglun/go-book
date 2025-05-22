package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// DoubanBook 豆瓣图书信息结构
type DoubanBook struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ISBN        string  `json:"isbn"`
	Rating      float64 `json:"rating"`
	Cover       string  `json:"cover"`
}

// DoubanService 豆瓣图书服务
type DoubanService struct {
	client *http.Client
}

// NewDoubanService 创建豆瓣服务实例
func NewDoubanService() *DoubanService {
	return &DoubanService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchBooks 搜索豆瓣图书
func (s *DoubanService) SearchBooks(keyword string) ([]DoubanBook, error) {
	// URL编码关键词
	searchURL := fmt.Sprintf("https://book.douban.com/subject_search?search_text=%s", url.QueryEscape(keyword))

	// 创建请求
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头，模拟浏览器行为
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var books []DoubanBook

	// 解析搜索结果
	doc.Find(".subject-item").Each(func(i int, s *goquery.Selection) {
		book := DoubanBook{}

		// 获取标题
		book.Title = strings.TrimSpace(s.Find(".title a").Text())

		// 获取作者
		book.Author = strings.TrimSpace(s.Find(".info .author").Text())

		// 获取描述
		book.Description = strings.TrimSpace(s.Find(".info p").Text())

		// 获取评分
		ratingStr := s.Find(".rating_nums").Text()
		if ratingStr != "" {
			fmt.Sscanf(ratingStr, "%f", &book.Rating)
		}

		// 获取封面图片
		book.Cover, _ = s.Find(".pic img").Attr("src")

		// 获取价格和ISBN（从详情页获取）
		if href, exists := s.Find(".title a").Attr("href"); exists {
			book.Price, book.ISBN = s.fetchBookDetail(href)
		}

		books = append(books, book)
	})

	return books, nil
}

// fetchBookDetail 获取图书详细信息
func (s *DoubanService) fetchBookDetail(url string) (float64, string) {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, ""
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, ""
	}

	// 获取价格
	var price float64
	priceStr := doc.Find(".subject-info span:contains('定价:')").Text()
	if priceStr != "" {
		priceStr = strings.TrimPrefix(priceStr, "定价:")
		priceStr = strings.TrimPrefix(priceStr, "CNY")
		priceStr = strings.TrimSpace(priceStr)
		fmt.Sscanf(priceStr, "%f", &price)
	}

	// 获取ISBN
	isbn := ""
	doc.Find(".subject-info span:contains('ISBN:')").Each(func(i int, s *goquery.Selection) {
		isbnText := s.Text()
		if strings.Contains(isbnText, "ISBN:") {
			isbn = strings.TrimSpace(strings.TrimPrefix(isbnText, "ISBN:"))
		}
	})

	return price, isbn
}