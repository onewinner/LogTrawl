package models

import "time"

// LogFile represents a log file with metadata
type LogFile struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	Lines        int       `json:"lines,omitempty"`        // 文件行数（可选）
	LastModified time.Time `json:"lastModified"`
	IsOpen       bool      `json:"isOpen"`
}

// LogContent represents the content of a log file
type LogContent struct {
	Lines []string `json:"lines"`
	Total int      `json:"total"`
}

// SearchResult represents search results
type SearchResult struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
	Matches    []int  `json:"matches"`
}

// RecentFile represents a recently opened file
type RecentFile struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	LastOpened   time.Time `json:"lastOpened"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}
