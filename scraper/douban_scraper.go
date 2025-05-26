package scraper

import (
	// "encoding/json"
	// "log"
	// "net/http"

	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	ISBN        string `json:"isbn"`
	Rating      string `json:"rating"`
	Cover       string `json:"cover"`
}

func DoubanScraper() []Book {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("douban.com", "book.douban.com"),
		colly.Async(true),
		colly.MaxDepth(2),
	)

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	// 设置并发限制
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*douban.*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		log.Println("Link found:", link)
		// Filter URLs that match the pattern https://book.douban.com/subject/[0-9]+
		if matched, _ := regexp.MatchString(`https://book\.douban\.com/subject/[0-9]+`, link); matched {
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			c.Visit(e.Request.AbsoluteURL(link))
		}
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	var books []Book

	c.OnHTML("div[class=info]", func(e *colly.HTMLElement) {
		title := e.ChildText("a[class=title]")
		author := e.ChildText("div[class=pub]")
		description := e.ChildText("p[class=abstract]")
		isbn := e.ChildText("span[class=isbn]")
		rating := e.ChildText("span[class=rating_nums]")
		cover := e.ChildAttr("a[class=nbg]", "href")
		book := Book{
			Title:       title,
			Author:      author,
			Description: description,
			ISBN:        isbn,
			Rating:      rating,
			Cover:       cover,
		}

		books = append(books, book)
	})

	// 添加错误处理
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if err := c.Visit("https://book.douban.com/subject_search?search_text=%E4%B8%89%E4%BD%93"); err != nil {
		log.Printf("Error visiting URL: %v", err)
	}

	c.Wait()

	// Return empty slice since no books were collected yet
	return []Book{}
}
