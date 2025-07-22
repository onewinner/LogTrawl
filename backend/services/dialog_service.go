package services

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DialogService struct {
	ctx context.Context
}

func NewDialogService(ctx context.Context) *DialogService {
	return &DialogService{
		ctx: ctx,
	}
}

// OpenFileDialog opens a file dialog and returns the selected file path
func (ds *DialogService) OpenFileDialog() (string, error) {
	filePath, err := runtime.OpenFileDialog(ds.ctx, runtime.OpenDialogOptions{
		Title: "选择日志文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "日志文件",
				Pattern:     "*.log;*.txt;*.out",
			},
			{
				DisplayName: "所有文件",
				Pattern:     "*.*",
			},
		},
	})
	return filePath, err
}

// OpenDirectoryDialog opens a directory dialog and returns the selected directory path
func (ds *DialogService) OpenDirectoryDialog() (string, error) {
	dirPath, err := runtime.OpenDirectoryDialog(ds.ctx, runtime.OpenDialogOptions{
		Title: "选择日志文件夹",
	})
	return dirPath, err
}

// SaveFileDialog opens a save file dialog
func (ds *DialogService) SaveFileDialog() (string, error) {
	filePath, err := runtime.SaveFileDialog(ds.ctx, runtime.SaveDialogOptions{
		Title:           "保存日志文件",
		DefaultFilename: "exported_logs.txt",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文本文件",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "日志文件",
				Pattern:     "*.log",
			},
			{
				DisplayName: "所有文件",
				Pattern:     "*.*",
			},
		},
	})
	return filePath, err
}
