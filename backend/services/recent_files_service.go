package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"LogTrawl/backend/models"
)

type RecentFilesService struct {
	configPath string
	maxFiles   int
}

func NewRecentFilesService() *RecentFilesService {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		// 如果获取失败，使用当前目录
		configDir = "."
	}
	
	appConfigDir := filepath.Join(configDir, "LogTrawl")
	os.MkdirAll(appConfigDir, 0755)
	
	return &RecentFilesService{
		configPath: filepath.Join(appConfigDir, "recent_files.json"),
		maxFiles:   10,
	}
}

// GetRecentFiles returns a list of recently opened files
func (rfs *RecentFilesService) GetRecentFiles() ([]models.LogFile, error) {
	recentFiles, err := rfs.loadRecentFiles()
	if err != nil {
		return []models.LogFile{}, nil // 如果加载失败，返回空列表
	}

	var logFiles []models.LogFile
	for _, rf := range recentFiles {
		// 检查文件是否仍然存在
		if _, err := os.Stat(rf.Path); err == nil {
			logFiles = append(logFiles, models.LogFile{
				ID:           fmt.Sprintf("recent_%d", time.Now().UnixNano()),
				Name:         rf.Name,
				Path:         rf.Path,
				Size:         rf.Size,
				LastModified: rf.LastModified,
				IsOpen:       false,
			})
		}
	}

	return logFiles, nil
}

// AddRecentFile adds a file to the recent files list
func (rfs *RecentFilesService) AddRecentFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	recentFiles, _ := rfs.loadRecentFiles()

	// 移除已存在的相同文件
	for i, rf := range recentFiles {
		if rf.Path == filePath {
			recentFiles = append(recentFiles[:i], recentFiles[i+1:]...)
			break
		}
	}

	// 添加新文件到列表开头
	newFile := models.RecentFile{
		Path:         filePath,
		Name:         fileInfo.Name(),
		LastOpened:   time.Now(),
		Size:         fileInfo.Size(),
		LastModified: fileInfo.ModTime(),
	}

	recentFiles = append([]models.RecentFile{newFile}, recentFiles...)

	// 限制最大文件数量
	if len(recentFiles) > rfs.maxFiles {
		recentFiles = recentFiles[:rfs.maxFiles]
	}

	return rfs.saveRecentFiles(recentFiles)
}

// RemoveRecentFile removes a file from the recent files list
func (rfs *RecentFilesService) RemoveRecentFile(filePath string) error {
	recentFiles, err := rfs.loadRecentFiles()
	if err != nil {
		return err
	}

	for i, rf := range recentFiles {
		if rf.Path == filePath {
			recentFiles = append(recentFiles[:i], recentFiles[i+1:]...)
			break
		}
	}

	return rfs.saveRecentFiles(recentFiles)
}

// ClearRecentFiles clears all recent files
func (rfs *RecentFilesService) ClearRecentFiles() error {
	return rfs.saveRecentFiles([]models.RecentFile{})
}

func (rfs *RecentFilesService) loadRecentFiles() ([]models.RecentFile, error) {
	data, err := os.ReadFile(rfs.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.RecentFile{}, nil
		}
		return nil, err
	}

	var recentFiles []models.RecentFile
	err = json.Unmarshal(data, &recentFiles)
	if err != nil {
		return []models.RecentFile{}, nil // 如果解析失败，返回空列表
	}

	// 按最后打开时间排序
	sort.Slice(recentFiles, func(i, j int) bool {
		return recentFiles[i].LastOpened.After(recentFiles[j].LastOpened)
	})

	return recentFiles, nil
}

func (rfs *RecentFilesService) saveRecentFiles(recentFiles []models.RecentFile) error {
	data, err := json.MarshalIndent(recentFiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(rfs.configPath, data, 0644)
}
