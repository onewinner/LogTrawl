package models

import "time"

// AnalysisResult 分析结果
type AnalysisResult struct {
	Overview    OverviewStats `json:"overview"`
	IPAnalysis  IPAnalysis    `json:"ipAnalysis"`
	URLAnalysis URLAnalysis   `json:"urlAnalysis"`
}

// OverviewStats 概览统计
type OverviewStats struct {
	TotalRequests int `json:"totalRequests"`
	UniqueIPs     int `json:"uniqueIPs"`
	InternalIPs   int `json:"internalIPs"`
	ExternalIPs   int `json:"externalIPs"`
}

// IPAnalysis IP分析结果
type IPAnalysis struct {
	InternalTop10 []IPStat `json:"internalTop10"`
	ExternalTop10 []IPStat `json:"externalTop10"`
}

// IPStat IP统计信息
type IPStat struct {
	Rank        int    `json:"rank"`
	IP          string `json:"ip"`
	Count       int    `json:"count"`
	FirstAccess string `json:"firstAccess"`
	LastAccess  string `json:"lastAccess"`
}

// URLAnalysis URL分析结果
type URLAnalysis struct {
	GetTop10  []URLStat `json:"getTop10"`
	PostTop10 []URLStat `json:"postTop10"`
}

// URLStat URL统计信息
type URLStat struct {
	Rank       int     `json:"rank"`
	URL        string  `json:"url"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// SpecificIPResult 特定IP分析结果
type SpecificIPResult struct {
	IP            string    `json:"ip"`
	IPType        string    `json:"ipType"`
	TotalRequests int       `json:"totalRequests"`
	FirstAccess   string    `json:"firstAccess"`
	LastAccess    string    `json:"lastAccess"`
	URLAnalysis   SpecificIPURLAnalysis `json:"urlAnalysis"`
}

// SpecificIPURLAnalysis 特定IP的URL分析
type SpecificIPURLAnalysis struct {
	GetTop10  []URLStat `json:"getTop10"`
	PostTop10 []URLStat `json:"postTop10"`
}

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time
	IP        string
	Method    string
	URL       string
	Status    int
	Size      int
	UserAgent string
	Referer   string
}

// IPInfo IP信息
type IPInfo struct {
	IP          string
	Count       int
	FirstAccess time.Time
	LastAccess  time.Time
	IsInternal  bool
}

// URLInfo URL信息
type URLInfo struct {
	URL    string
	Method string
	Count  int
}

// DatePattern 日期模式
type DatePattern struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Example     string `json:"example"`
	Description string `json:"description"`
}

// DateRange 文件日期范围
type DateRange struct {
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	DatePattern  string     `json:"datePattern"`
	TotalDays    float64    `json:"totalDays"`
	HasDateInfo  bool       `json:"hasDateInfo"`
}

// FileSplitOptions 文件分片选项
type FileSplitOptions struct {
	Strategy     string  `json:"strategy"`
	FilePath     string  `json:"filePath"`
	OutputDir    string  `json:"outputDir"`
	DatePattern  string  `json:"datePattern,omitempty"`
	DaysPerFile  float64 `json:"daysPerFile,omitempty"`  // 改为float64支持小数
	SizePerFile  int64   `json:"sizePerFile,omitempty"`
	LinesPerFile int     `json:"linesPerFile,omitempty"`
}

// FileSplitResult 文件分片结果
type FileSplitResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	OutputFiles []string `json:"outputFiles"`
	TotalFiles  int      `json:"totalFiles"`
	TotalSize   int64    `json:"totalSize"`
}
