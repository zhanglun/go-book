package scraper

import (
	// "encoding/json"
	// "log"
	// "net/http"

	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

type Book struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ISBN        string  `json:"isbn"`
	Rating      float64 `json:"rating"`
	Cover       string  `json:"cover"`
}

func DoubanScraper() []Book {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("douban.com", "book.douban.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		log.Println("Link found:", link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if err := c.Visit("https://book.douban.com/subject_search?search_text=%E4%B8%89%E4%BD%93"); err != nil {
		log.Printf("Error visiting URL: %v", err)
	}
	// Return empty slice since no books were collected yet
	return []Book{}
}
