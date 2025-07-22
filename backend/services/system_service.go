package services

import (
	"os"
	"runtime"
	"time"

	"LogTrawl/backend/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemService 系统资源监控服务
type SystemService struct {
	lastCPUTime  time.Time
	lastCPUUsage float64
	process      *process.Process
}

// NewSystemService 创建新的系统服务
func NewSystemService() *SystemService {
	// 获取当前进程
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		proc = nil // 如果获取失败，设为nil
	}

	return &SystemService{
		lastCPUTime: time.Now(),
		process:     proc,
	}
}

// GetSystemInfo 获取系统资源信息
func (s *SystemService) GetSystemInfo() (*models.SystemInfo, error) {
	// 获取应用程序CPU使用率
	appCPUUsage, err := s.getAppCPUUsage()
	if err != nil {
		appCPUUsage = 0 // 如果获取失败，设为0
	}

	// 获取系统内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取应用程序内存使用
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 计算应用程序内存占系统总内存的百分比
	appMemoryPercent := float64(memStats.Alloc) / float64(memInfo.Total) * 100

	systemInfo := &models.SystemInfo{
		CPUUsage:        appCPUUsage,
		MemoryUsage:     appMemoryPercent,
		MemoryUsed:      memInfo.Used,
		MemoryTotal:     memInfo.Total,
		MemoryAvailable: memInfo.Available,
		AppMemoryUsage:  memStats.Alloc,
	}

	return systemInfo, nil
}

// getAppCPUUsage 获取应用程序CPU使用率
func (s *SystemService) getAppCPUUsage() (float64, error) {
	if s.process == nil {
		return 0, nil
	}

	// 获取进程CPU使用率
	cpuPercent, err := s.process.CPUPercent()
	if err != nil {
		return 0, err
	}

	return cpuPercent, nil
}

// getCPUUsage 获取系统整体CPU使用率（备用方法）
func (s *SystemService) getCPUUsage() (float64, error) {
	// 使用gopsutil获取CPU使用率
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}

	if len(percentages) > 0 {
		s.lastCPUUsage = percentages[0]
		s.lastCPUTime = time.Now()
		return percentages[0], nil
	}

	return s.lastCPUUsage, nil
}
