package scraper

import (
	"testing"
)

func TestFindList(t *testing.T) {
	result := DoubanScraper()
	t.Logf("Scraper result: %v", result)
}
