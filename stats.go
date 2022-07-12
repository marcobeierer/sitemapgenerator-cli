package main

import "time"

// NOTE copy in sitemapgenerator/data/stats.go
type Stats struct {
	URLLimitReached            bool
	CrawledResourcesCount      int
	DeadResourcesCount         int
	TimedOutResourcesCount     int
	SitemapURLCount            int64
	SitemapImageCount          int
	SitemapVideoCount          int
	StartedAt                  time.Time
	FinishedAt                 time.Time
	TokenUsed                  bool
	MaxFetchers                int
	URLLimit                   int
	SitemapIndexCount          int // number of available subsitemaps
	SitemapIndexNumberOfDigits int // if for example 6, the value is 000001
	CrawlDelayInSeconds        float64
	// SitemapIndexFilenames  []string
	// SitemapIndexURLs           []*SitemapIndexURL `json:",omitempty"`
}
