package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"LogTrawl/backend/models"
)

type FileSplitterService struct{}

func NewFileSplitterService() *FileSplitterService {
	return &FileSplitterService{}
}

// SplitStrategy 分片策略
type SplitStrategy string

const (
	SplitByDate SplitStrategy = "date"
	SplitBySize SplitStrategy = "size"
	SplitByLines SplitStrategy = "lines"
)

// SplitOptions 分片选项
type SplitOptions struct {
	Strategy     SplitStrategy `json:"strategy"`
	FilePath     string        `json:"filePath"`
	OutputDir    string        `json:"outputDir"`

	// 按日期分片
	DatePattern  string        `json:"datePattern,omitempty"`  // 日期正则表达式
	DaysPerFile  float64       `json:"daysPerFile,omitempty"`  // 每个文件包含的天数（支持小数）

	// 按大小分片
	SizePerFile  int64         `json:"sizePerFile,omitempty"`  // 每个文件的大小（字节）

	// 按行数分片
	LinesPerFile int           `json:"linesPerFile,omitempty"` // 每个文件的行数
}

// SplitResult 分片结果
type SplitResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	OutputFiles []string `json:"outputFiles"`
	TotalFiles  int      `json:"totalFiles"`
	TotalSize   int64    `json:"totalSize"`
}

// SplitFile 分片文件
func (fs *FileSplitterService) SplitFile(options SplitOptions) (*SplitResult, error) {
	// 检查输入文件
	if _, err := os.Stat(options.FilePath); os.IsNotExist(err) {
		return &SplitResult{
			Success: false,
			Message: "输入文件不存在",
		}, err
	}

	// 创建输出目录
	if err := os.MkdirAll(options.OutputDir, 0755); err != nil {
		return &SplitResult{
			Success: false,
			Message: "无法创建输出目录",
		}, err
	}

	switch options.Strategy {
	case SplitByDate:
		return fs.splitByDate(options)
	case SplitBySize:
		return fs.splitBySize(options)
	case SplitByLines:
		return fs.splitByLines(options)
	default:
		return &SplitResult{
			Success: false,
			Message: "不支持的分片策略",
		}, fmt.Errorf("unsupported split strategy: %s", options.Strategy)
	}
}

// splitByDate 按日期分片
func (fs *FileSplitterService) splitByDate(options SplitOptions) (*SplitResult, error) {
	file, err := os.Open(options.FilePath)
	if err != nil {
		return &SplitResult{Success: false, Message: "无法打开文件"}, err
	}
	defer file.Close()

	// 获取原文件名（不含扩展名）
	originalName := filepath.Base(options.FilePath)
	ext := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, ext)

	// 自动检测日期格式
	datePatterns := []struct {
		pattern *regexp.Regexp
		layout  string
	}{
		{regexp.MustCompile(`(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})`), "2006-01-02 15:04:05"},
		{regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`), "2006-01-02"},
		{regexp.MustCompile(`(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2})`), "02/Jan/2006:15:04:05"},
		{regexp.MustCompile(`(\d{2}/\d{2}/\d{4}\s+\d{2}:\d{2}:\d{2})`), "01/02/2006 15:04:05"},
	}

	var dateRegex *regexp.Regexp
	var timeLayout string

	// 检测日期格式
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && dateRegex == nil {
		line := scanner.Text()
		for _, dp := range datePatterns {
			if dp.pattern.MatchString(line) {
				dateRegex = dp.pattern
				timeLayout = dp.layout
				break
			}
		}
	}

	if dateRegex == nil {
		return &SplitResult{Success: false, Message: "无法检测到日期格式"}, fmt.Errorf("no date pattern found")
	}

	// 重新开始读取文件
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024 // 1MB缓冲区
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var outputFiles []string
	var currentFile *os.File
	var currentWriter *bufio.Writer
	var currentStartDate, currentEndDate time.Time
	var totalSize int64
	var hasCurrentFile bool

	// 计算分片时长（小时）
	hoursPerFile := options.DaysPerFile * 24

	for scanner.Scan() {
		line := scanner.Text()

		// 提取日期
		matches := dateRegex.FindStringSubmatch(line)
		if len(matches) <= 1 {
			// 如果没有日期，写入当前文件（如果有的话）
			if currentWriter != nil {
				lineBytes := []byte(line + "\n")
				currentWriter.Write(lineBytes)
				totalSize += int64(len(lineBytes))
			}
			continue
		}

		// 解析日期
		lineDate, err := time.Parse(timeLayout, matches[1])
		if err != nil {
			// 日期解析失败，跳过这行或写入当前文件
			if currentWriter != nil {
				lineBytes := []byte(line + "\n")
				currentWriter.Write(lineBytes)
				totalSize += int64(len(lineBytes))
			}
			continue
		}

		// 检查是否需要创建新文件
		needNewFile := false
		if !hasCurrentFile {
			needNewFile = true
			currentStartDate = lineDate
		} else {
			// 检查是否超过了时间范围
			duration := lineDate.Sub(currentStartDate)
			if duration.Hours() >= hoursPerFile {
				needNewFile = true
			}
		}

		if needNewFile {
			// 关闭当前文件
			if currentWriter != nil {
				currentWriter.Flush()
				currentFile.Close()
			}

			// 设置新的时间范围
			currentStartDate = lineDate
			currentEndDate = currentStartDate.Add(time.Duration(hoursPerFile * float64(time.Hour)))

			// 创建新文件名：原文件名_开始日期_结束日期.扩展名
			startStr := currentStartDate.Format("2006-01-02")
			endStr := currentEndDate.Format("2006-01-02")
			fileName := fmt.Sprintf("%s_%s_%s%s", baseName, startStr, endStr, ext)
			filePath := filepath.Join(options.OutputDir, fileName)

			currentFile, err = os.Create(filePath)
			if err != nil {
				return &SplitResult{Success: false, Message: "无法创建输出文件"}, err
			}

			currentWriter = bufio.NewWriter(currentFile)
			outputFiles = append(outputFiles, filePath)
			hasCurrentFile = true
		}

		// 写入当前行
		if currentWriter != nil {
			lineBytes := []byte(line + "\n")
			currentWriter.Write(lineBytes)
			totalSize += int64(len(lineBytes))
		}
	}

	// 关闭最后一个文件
	if currentWriter != nil {
		currentWriter.Flush()
		currentFile.Close()
	}

	if err := scanner.Err(); err != nil {
		return &SplitResult{Success: false, Message: "读取文件失败"}, err
	}

	return &SplitResult{
		Success:     true,
		Message:     fmt.Sprintf("成功按日期分片为 %d 个文件", len(outputFiles)),
		OutputFiles: outputFiles,
		TotalFiles:  len(outputFiles),
		TotalSize:   totalSize,
	}, nil
}

// splitBySize 按大小分片
func (fs *FileSplitterService) splitBySize(options SplitOptions) (*SplitResult, error) {
	file, err := os.Open(options.FilePath)
	if err != nil {
		return &SplitResult{Success: false, Message: "无法打开文件"}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024 // 1MB缓冲区
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var outputFiles []string
	var currentFile *os.File
	var currentWriter *bufio.Writer
	var currentSize int64
	var fileIndex int
	var totalSize int64

	baseFileName := strings.TrimSuffix(filepath.Base(options.FilePath), filepath.Ext(options.FilePath))

	createNewFile := func() error {
		if currentWriter != nil {
			currentWriter.Flush()
			currentFile.Close()
		}

		fileIndex++
		fileName := fmt.Sprintf("%s_part%d.log", baseFileName, fileIndex)
		filePath := filepath.Join(options.OutputDir, fileName)
		
		var err error
		currentFile, err = os.Create(filePath)
		if err != nil {
			return err
		}
		
		currentWriter = bufio.NewWriter(currentFile)
		outputFiles = append(outputFiles, filePath)
		currentSize = 0
		return nil
	}

	// 创建第一个文件
	if err := createNewFile(); err != nil {
		return &SplitResult{Success: false, Message: "无法创建输出文件"}, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := []byte(line + "\n")
		lineSize := int64(len(lineBytes))

		// 检查是否需要创建新文件
		if currentSize+lineSize > options.SizePerFile && currentSize > 0 {
			if err := createNewFile(); err != nil {
				return &SplitResult{Success: false, Message: "无法创建输出文件"}, err
			}
		}

		// 写入当前行
		currentWriter.Write(lineBytes)
		currentSize += lineSize
		totalSize += lineSize
	}

	// 关闭最后一个文件
	if currentWriter != nil {
		currentWriter.Flush()
		currentFile.Close()
	}

	if err := scanner.Err(); err != nil {
		return &SplitResult{Success: false, Message: "读取文件失败"}, err
	}

	return &SplitResult{
		Success:     true,
		Message:     fmt.Sprintf("成功按大小分片为 %d 个文件", len(outputFiles)),
		OutputFiles: outputFiles,
		TotalFiles:  len(outputFiles),
		TotalSize:   totalSize,
	}, nil
}

// splitByLines 按行数分片
func (fs *FileSplitterService) splitByLines(options SplitOptions) (*SplitResult, error) {
	file, err := os.Open(options.FilePath)
	if err != nil {
		return &SplitResult{Success: false, Message: "无法打开文件"}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxCapacity = 1024 * 1024 // 1MB缓冲区
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var outputFiles []string
	var currentFile *os.File
	var currentWriter *bufio.Writer
	var currentLines int
	var fileIndex int
	var totalSize int64

	baseFileName := strings.TrimSuffix(filepath.Base(options.FilePath), filepath.Ext(options.FilePath))

	createNewFile := func() error {
		if currentWriter != nil {
			currentWriter.Flush()
			currentFile.Close()
		}

		fileIndex++
		fileName := fmt.Sprintf("%s_part%d.log", baseFileName, fileIndex)
		filePath := filepath.Join(options.OutputDir, fileName)
		
		var err error
		currentFile, err = os.Create(filePath)
		if err != nil {
			return err
		}
		
		currentWriter = bufio.NewWriter(currentFile)
		outputFiles = append(outputFiles, filePath)
		currentLines = 0
		return nil
	}

	// 创建第一个文件
	if err := createNewFile(); err != nil {
		return &SplitResult{Success: false, Message: "无法创建输出文件"}, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := []byte(line + "\n")

		// 检查是否需要创建新文件
		if currentLines >= options.LinesPerFile && currentLines > 0 {
			if err := createNewFile(); err != nil {
				return &SplitResult{Success: false, Message: "无法创建输出文件"}, err
			}
		}

		// 写入当前行
		currentWriter.Write(lineBytes)
		currentLines++
		totalSize += int64(len(lineBytes))
	}

	// 关闭最后一个文件
	if currentWriter != nil {
		currentWriter.Flush()
		currentFile.Close()
	}

	if err := scanner.Err(); err != nil {
		return &SplitResult{Success: false, Message: "读取文件失败"}, err
	}

	return &SplitResult{
		Success:     true,
		Message:     fmt.Sprintf("成功按行数分片为 %d 个文件", len(outputFiles)),
		OutputFiles: outputFiles,
		TotalFiles:  len(outputFiles),
		TotalSize:   totalSize,
	}, nil
}

// GetCommonDatePatterns 获取常用日期模式
func (fs *FileSplitterService) GetCommonDatePatterns() []models.DatePattern {
	return []models.DatePattern{
		{
			Name:        "标准日期 (YYYY-MM-DD)",
			Pattern:     `(\d{4}-\d{2}-\d{2})`,
			Example:     "2023-12-25",
			Description: "匹配 YYYY-MM-DD 格式的日期",
		},
		{
			Name:        "美式日期 (MM/DD/YYYY)",
			Pattern:     `(\d{2}/\d{2}/\d{4})`,
			Example:     "12/25/2023",
			Description: "匹配 MM/DD/YYYY 格式的日期",
		},
		{
			Name:        "欧式日期 (DD/MM/YYYY)",
			Pattern:     `(\d{2}/\d{2}/\d{4})`,
			Example:     "25/12/2023",
			Description: "匹配 DD/MM/YYYY 格式的日期",
		},
		{
			Name:        "Apache日志日期",
			Pattern:     `\[(\d{2}/\w{3}/\d{4})`,
			Example:     "[25/Dec/2023",
			Description: "匹配 Apache 日志中的日期格式",
		},
		{
			Name:        "时间戳日期 (YYYY-MM-DD HH:MM:SS)",
			Pattern:     `(\d{4}-\d{2}-\d{2})\s+\d{2}:\d{2}:\d{2}`,
			Example:     "2023-12-25 10:30:45",
			Description: "匹配完整时间戳中的日期部分",
		},
	}
}
