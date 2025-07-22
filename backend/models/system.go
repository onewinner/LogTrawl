package models

// SystemInfo 系统资源信息
type SystemInfo struct {
	// CPU 使用率 (百分比)
	CPUUsage float64 `json:"cpuUsage"`
	
	// 内存信息
	MemoryUsage     float64 `json:"memoryUsage"`     // 内存使用率 (百分比)
	MemoryUsed      uint64  `json:"memoryUsed"`      // 已使用内存 (字节)
	MemoryTotal     uint64  `json:"memoryTotal"`     // 总内存 (字节)
	MemoryAvailable uint64  `json:"memoryAvailable"` // 可用内存 (字节)
	
	// 应用程序内存使用
	AppMemoryUsage uint64 `json:"appMemoryUsage"` // 应用程序内存使用 (字节)
}
