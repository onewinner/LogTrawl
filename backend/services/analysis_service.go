package services

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"LogTrawl/backend/models"
)

type AnalysisService struct{
	mu sync.RWMutex
	progressCallbacks map[string]func(progress int, status string)
	cancelChannels map[string]chan bool
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(progress int, status string)

func NewAnalysisService() *AnalysisService {
	return &AnalysisService{
		progressCallbacks: make(map[string]func(progress int, status string)),
		cancelChannels: make(map[string]chan bool),
	}
}

// AnalyzeLogFile 分析日志文件
func (as *AnalysisService) AnalyzeLogFile(filePath string) (*models.AnalysisResult, error) {
	return as.AnalyzeLogFileWithProgress(filePath, "")
}

// AnalyzeLogFileWithProgress 分析日志文件（支持进度回调）
func (as *AnalysisService) AnalyzeLogFileWithProgress(filePath, sessionID string) (*models.AnalysisResult, error) {
	fmt.Printf("开始分析日志文件: %s\n", filePath)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法获取文件信息: %v", err)
	}

	fileSize := fileInfo.Size()
	maxAnalysisSize := int64(1024 * 1024 * 1024) // 1GB最大分析限制

	// 检查文件大小是否超过1GB限制
	if fileSize > maxAnalysisSize {
		return nil, fmt.Errorf("文件大小 %.2f GB 超过最大分析限制 1GB，请使用文件分片功能", float64(fileSize)/(1024*1024*1024))
	}

	isLargeFile := fileSize > 200*1024*1024 // 200MB以上认为是大文件

	if isLargeFile {
		fmt.Printf("检测到大文件 (%.2f MB)，使用优化分析模式\n", float64(fileSize)/(1024*1024))
		return as.analyzeLargeFile(filePath, sessionID, fileSize)
	}

	return as.analyzeSmallFile(filePath, sessionID)
}

// analyzeSmallFile 分析小文件（原有逻辑）
func (as *AnalysisService) analyzeSmallFile(filePath, sessionID string) (*models.AnalysisResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 存储分析数据
	ipMap := make(map[string]*models.IPInfo)
	urlMap := make(map[string]*models.URLInfo)
	totalRequests := 0
	parsedLines := 0

	// 日志解析正则表达式
	logRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+\S+\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+[^"]*"\s+(\d+)\s+(\d+)`)

	scanner := bufio.NewScanner(file)
	// 使用更大的缓冲区
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parsedLines++

		// 进度回调
		if sessionID != "" && parsedLines%10000 == 0 {
			as.updateProgress(sessionID, 50, fmt.Sprintf("已处理 %d 行", parsedLines))
		}

		if line == "" {
			continue
		}

		entry := as.parseLogLine(line, logRegex)
		if entry == nil {
			continue
		}

		totalRequests++

		// 处理IP统计
		as.processIPStats(entry, ipMap)

		// 处理URL统计
		as.processURLStats(entry, urlMap)
	}

	if sessionID != "" {
		as.updateProgress(sessionID, 80, "正在生成分析结果...")
	}

	fmt.Printf("分析完成: 总行数=%d, 有效请求=%d, 独立IP=%d\n", parsedLines, totalRequests, len(ipMap))

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 生成分析结果
	result := &models.AnalysisResult{
		Overview:    as.generateOverview(totalRequests, ipMap),
		IPAnalysis:  as.generateIPAnalysis(ipMap),
		URLAnalysis: as.generateURLAnalysis(urlMap, totalRequests),
	}

	// 确保所有切片都已初始化
	if result.IPAnalysis.InternalTop10 == nil {
		result.IPAnalysis.InternalTop10 = []models.IPStat{}
	}
	if result.IPAnalysis.ExternalTop10 == nil {
		result.IPAnalysis.ExternalTop10 = []models.IPStat{}
	}
	if result.URLAnalysis.GetTop10 == nil {
		result.URLAnalysis.GetTop10 = []models.URLStat{}
	}
	if result.URLAnalysis.PostTop10 == nil {
		result.URLAnalysis.PostTop10 = []models.URLStat{}
	}

	if sessionID != "" {
		as.updateProgress(sessionID, 100, "分析完成")
	}

	return result, nil
}

// analyzeLargeFile 分析大文件（分块处理）
func (as *AnalysisService) analyzeLargeFile(filePath, sessionID string, fileSize int64) (*models.AnalysisResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 存储分析数据
	ipMap := make(map[string]*models.IPInfo)
	urlMap := make(map[string]*models.URLInfo)
	totalRequests := 0
	parsedLines := 0
	processedBytes := int64(0)

	// 使用更大的缓冲区处理大文件
	const bufferSize = 2 * 1024 * 1024 // 2MB缓冲区
	scanner := bufio.NewScanner(file)
	buf := make([]byte, bufferSize)
	scanner.Buffer(buf, bufferSize)

	// 日志解析正则表达式
	logRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+\S+\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+[^"]*"\s+(\d+)\s+(\d+)`)

	// 检查是否需要取消
	cancelChan := as.getCancelChannel(sessionID)

	for scanner.Scan() {
		// 检查取消信号
		select {
		case <-cancelChan:
			return nil, fmt.Errorf("分析已取消")
		default:
		}

		line := strings.TrimSpace(scanner.Text())
		parsedLines++
		processedBytes += int64(len(line) + 1) // +1 for newline

		// 更频繁的进度更新（每5000行）
		if sessionID != "" && parsedLines%5000 == 0 {
			progress := int(float64(processedBytes) / float64(fileSize) * 80) // 80%用于处理，20%用于结果生成
			as.updateProgress(sessionID, progress, fmt.Sprintf("已处理 %d 行 (%.1f%%)", parsedLines, float64(processedBytes)/float64(fileSize)*100))
		}

		if line == "" {
			continue
		}

		entry := as.parseLogLine(line, logRegex)
		if entry == nil {
			continue
		}

		totalRequests++

		// 处理IP统计
		as.processIPStats(entry, ipMap)

		// 处理URL统计
		as.processURLStats(entry, urlMap)

		// 定期清理内存（每100万行）
		if parsedLines%1000000 == 0 {
			// 强制垃圾回收
			// runtime.GC() // 可选：强制垃圾回收
		}
	}

	if sessionID != "" {
		as.updateProgress(sessionID, 85, "正在生成分析结果...")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	fmt.Printf("大文件分析完成: 总行数=%d, 有效请求=%d, 独立IP=%d\n", parsedLines, totalRequests, len(ipMap))

	// 生成分析结果
	result := &models.AnalysisResult{
		Overview:    as.generateOverview(totalRequests, ipMap),
		IPAnalysis:  as.generateIPAnalysis(ipMap),
		URLAnalysis: as.generateURLAnalysis(urlMap, totalRequests),
	}

	// 确保所有切片都已初始化
	if result.IPAnalysis.InternalTop10 == nil {
		result.IPAnalysis.InternalTop10 = []models.IPStat{}
	}
	if result.IPAnalysis.ExternalTop10 == nil {
		result.IPAnalysis.ExternalTop10 = []models.IPStat{}
	}
	if result.URLAnalysis.GetTop10 == nil {
		result.URLAnalysis.GetTop10 = []models.URLStat{}
	}
	if result.URLAnalysis.PostTop10 == nil {
		result.URLAnalysis.PostTop10 = []models.URLStat{}
	}

	if sessionID != "" {
		as.updateProgress(sessionID, 100, "分析完成")
	}

	return result, nil
}

// updateProgress 更新进度
func (as *AnalysisService) updateProgress(sessionID string, progress int, status string) {
	as.mu.RLock()
	callback, exists := as.progressCallbacks[sessionID]
	as.mu.RUnlock()

	if exists && callback != nil {
		callback(progress, status)
	}
}

// getCancelChannel 获取取消通道
func (as *AnalysisService) getCancelChannel(sessionID string) <-chan bool {
	as.mu.RLock()
	defer as.mu.RUnlock()

	if ch, exists := as.cancelChannels[sessionID]; exists {
		return ch
	}

	// 创建一个永远不会发送信号的通道
	ch := make(chan bool)
	return ch
}

// SetProgressCallback 设置进度回调
func (as *AnalysisService) SetProgressCallback(sessionID string, callback ProgressCallback) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.progressCallbacks[sessionID] = callback
}

// CancelAnalysis 取消分析
func (as *AnalysisService) CancelAnalysis(sessionID string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if ch, exists := as.cancelChannels[sessionID]; exists {
		close(ch)
		delete(as.cancelChannels, sessionID)
	}
	delete(as.progressCallbacks, sessionID)
}

// StartAnalysis 开始分析（创建会话）
func (as *AnalysisService) StartAnalysis(sessionID string) {
	as.mu.Lock()
	defer as.mu.Unlock()
	as.cancelChannels[sessionID] = make(chan bool)
}

// parseLogLine 解析日志行
func (as *AnalysisService) parseLogLine(line string, regex *regexp.Regexp) *models.LogEntry {
	// 尝试多种日志格式
	entry := as.tryParseStandardLog(line)
	if entry != nil {
		return entry
	}

	// 如果标准格式失败，尝试简单解析
	return as.parseSimpleLogLine(line)
}

// tryParseStandardLog 尝试解析标准日志格式
func (as *AnalysisService) tryParseStandardLog(line string) *models.LogEntry {
	// 标准Apache/Nginx日志格式
	// 192.168.1.1 - - [25/Dec/2023:10:00:00 +0000] "GET /api/test HTTP/1.1" 200 1234
	standardRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+\S+\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+[^"]*"\s+(\d+)\s+(\d+)`)
	matches := standardRegex.FindStringSubmatch(line)
	if len(matches) >= 7 {
		status, _ := strconv.Atoi(matches[5])
		size, _ := strconv.Atoi(matches[6])
		return &models.LogEntry{
			Timestamp: time.Now(),
			IP:        matches[1],
			Method:    matches[3],
			URL:       matches[4],
			Status:    status,
			Size:      size,
		}
	}

	// 简化的日志格式
	// 192.168.1.1 "GET /api/test" 200 1234
	simpleRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+).*?"(\S+)\s+(\S+)".*?(\d+)\s+(\d+)`)
	matches = simpleRegex.FindStringSubmatch(line)
	if len(matches) >= 6 {
		status, _ := strconv.Atoi(matches[4])
		size, _ := strconv.Atoi(matches[5])
		return &models.LogEntry{
			Timestamp: time.Now(),
			IP:        matches[1],
			Method:    matches[2],
			URL:       matches[3],
			Status:    status,
			Size:      size,
		}
	}

	// 更简单的格式
	// 192.168.1.1 GET /api/test 200
	basicRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+(\S+)\s+(\S+)\s+(\d+)`)
	matches = basicRegex.FindStringSubmatch(line)
	if len(matches) >= 5 {
		status, _ := strconv.Atoi(matches[4])
		return &models.LogEntry{
			Timestamp: time.Now(),
			IP:        matches[1],
			Method:    matches[2],
			URL:       matches[3],
			Status:    status,
			Size:      0,
		}
	}

	return nil
}

// parseSimpleLogLine 简单解析日志行，提取IP和基本信息
func (as *AnalysisService) parseSimpleLogLine(line string) *models.LogEntry {
	// 尝试提取IP地址 - 更严格的IP格式
	ipRegex := regexp.MustCompile(`\b(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\b`)
	ipMatches := ipRegex.FindStringSubmatch(line)
	if len(ipMatches) < 2 {
		return nil
	}

	// 验证IP地址的有效性
	ip := ipMatches[1]
	if !as.isValidIP(ip) {
		return nil
	}

	// 尝试提取HTTP方法
	methodRegex := regexp.MustCompile(`\b(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH)\b`)
	methodMatches := methodRegex.FindStringSubmatch(line)
	method := "GET" // 默认方法
	if len(methodMatches) >= 2 {
		method = methodMatches[1]
	}

	// 尝试提取URL路径 - 改进的正则表达式
	urlRegex := regexp.MustCompile(`\b(GET|POST|PUT|DELETE|HEAD|OPTIONS|PATCH)\s+([^\s"]+)`)
	urlMatches := urlRegex.FindStringSubmatch(line)
	url := "/" // 默认URL
	if len(urlMatches) >= 3 {
		url = urlMatches[2]
		// 清理URL，移除可能的引号
		url = strings.Trim(url, `"'`)
	}

	return &models.LogEntry{
		Timestamp: time.Now(),
		IP:        ip,
		Method:    method,
		URL:       url,
		Status:    200, // 默认状态码
		Size:      0,   // 默认大小
	}
}

// isValidIP 验证IP地址的有效性
func (as *AnalysisService) isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// processIPStats 处理IP统计
func (as *AnalysisService) processIPStats(entry *models.LogEntry, ipMap map[string]*models.IPInfo) {
	if info, exists := ipMap[entry.IP]; exists {
		info.Count++
		if entry.Timestamp.After(info.LastAccess) {
			info.LastAccess = entry.Timestamp
		}
		if entry.Timestamp.Before(info.FirstAccess) {
			info.FirstAccess = entry.Timestamp
		}
	} else {
		ipMap[entry.IP] = &models.IPInfo{
			IP:          entry.IP,
			Count:       1,
			FirstAccess: entry.Timestamp,
			LastAccess:  entry.Timestamp,
			IsInternal:  as.isInternalIP(entry.IP),
		}
	}
}

// processURLStats 处理URL统计
func (as *AnalysisService) processURLStats(entry *models.LogEntry, urlMap map[string]*models.URLInfo) {
	key := fmt.Sprintf("%s:%s", entry.Method, entry.URL)
	if info, exists := urlMap[key]; exists {
		info.Count++
	} else {
		urlMap[key] = &models.URLInfo{
			URL:    entry.URL,
			Method: entry.Method,
			Count:  1,
		}
	}
}

// isInternalIP 判断是否为内网IP
func (as *AnalysisService) isInternalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 检查私有IP地址范围
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// generateOverview 生成概览统计
func (as *AnalysisService) generateOverview(totalRequests int, ipMap map[string]*models.IPInfo) models.OverviewStats {
	uniqueIPs := len(ipMap)
	internalIPs := 0
	externalIPs := 0

	for _, info := range ipMap {
		if info.IsInternal {
			internalIPs++
		} else {
			externalIPs++
		}
	}

	return models.OverviewStats{
		TotalRequests: totalRequests,
		UniqueIPs:     uniqueIPs,
		InternalIPs:   internalIPs,
		ExternalIPs:   externalIPs,
	}
}

// generateIPAnalysis 生成IP分析
func (as *AnalysisService) generateIPAnalysis(ipMap map[string]*models.IPInfo) models.IPAnalysis {
	internalIPs := make([]*models.IPInfo, 0)
	externalIPs := make([]*models.IPInfo, 0)

	for _, info := range ipMap {
		if info.IsInternal {
			internalIPs = append(internalIPs, info)
		} else {
			externalIPs = append(externalIPs, info)
		}
	}

	// 排序
	sort.Slice(internalIPs, func(i, j int) bool {
		return internalIPs[i].Count > internalIPs[j].Count
	})
	sort.Slice(externalIPs, func(i, j int) bool {
		return externalIPs[i].Count > externalIPs[j].Count
	})

	// 生成TOP10
	internalTop10 := as.generateIPTop10(internalIPs)
	externalTop10 := as.generateIPTop10(externalIPs)

	return models.IPAnalysis{
		InternalTop10: internalTop10,
		ExternalTop10: externalTop10,
	}
}

// generateIPTop10 生成IP TOP10
func (as *AnalysisService) generateIPTop10(ips []*models.IPInfo) []models.IPStat {
	result := make([]models.IPStat, 0) // 确保初始化为空切片而不是nil
	limit := 10
	if len(ips) < limit {
		limit = len(ips)
	}

	for i := 0; i < limit; i++ {
		info := ips[i]
		result = append(result, models.IPStat{
			Rank:        i + 1,
			IP:          info.IP,
			Count:       info.Count,
			FirstAccess: info.FirstAccess.Format("2006-01-02 15:04:05"),
			LastAccess:  info.LastAccess.Format("2006-01-02 15:04:05"),
		})
	}

	return result
}

// generateURLAnalysis 生成URL分析
func (as *AnalysisService) generateURLAnalysis(urlMap map[string]*models.URLInfo, totalRequests int) models.URLAnalysis {
	getURLs := make([]*models.URLInfo, 0)
	postURLs := make([]*models.URLInfo, 0)

	for _, info := range urlMap {
		if info.Method == "GET" {
			getURLs = append(getURLs, info)
		} else if info.Method == "POST" {
			postURLs = append(postURLs, info)
		}
	}

	// 排序
	sort.Slice(getURLs, func(i, j int) bool {
		return getURLs[i].Count > getURLs[j].Count
	})
	sort.Slice(postURLs, func(i, j int) bool {
		return postURLs[i].Count > postURLs[j].Count
	})

	// 生成TOP10
	getTop10 := as.generateURLTop10(getURLs, totalRequests)
	postTop10 := as.generateURLTop10(postURLs, totalRequests)

	return models.URLAnalysis{
		GetTop10:  getTop10,
		PostTop10: postTop10,
	}
}

// generateURLTop10 生成URL TOP10
func (as *AnalysisService) generateURLTop10(urls []*models.URLInfo, totalRequests int) []models.URLStat {
	result := make([]models.URLStat, 0) // 确保初始化为空切片而不是nil
	limit := 10
	if len(urls) < limit {
		limit = len(urls)
	}

	for i := 0; i < limit; i++ {
		info := urls[i]
		percentage := float64(info.Count) / float64(totalRequests) * 100
		result = append(result, models.URLStat{
			Rank:       i + 1,
			URL:        info.URL,
			Count:      info.Count,
			Percentage: percentage,
		})
	}

	return result
}

// AnalyzeSpecificIP 分析特定IP
func (as *AnalysisService) AnalyzeSpecificIP(filePath, targetIP string) (*models.SpecificIPResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	urlMap := make(map[string]*models.URLInfo)
	var firstAccess, lastAccess time.Time
	totalRequests := 0

	logRegex := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+)\s+\S+\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+[^"]*"\s+(\d+)\s+(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := as.parseLogLine(line, logRegex)
		if entry == nil || entry.IP != targetIP {
			continue
		}

		totalRequests++

		// 记录时间范围
		if firstAccess.IsZero() || entry.Timestamp.Before(firstAccess) {
			firstAccess = entry.Timestamp
		}
		if lastAccess.IsZero() || entry.Timestamp.After(lastAccess) {
			lastAccess = entry.Timestamp
		}

		// 统计URL
		key := fmt.Sprintf("%s:%s", entry.Method, entry.URL)
		if info, exists := urlMap[key]; exists {
			info.Count++
		} else {
			urlMap[key] = &models.URLInfo{
				URL:    entry.URL,
				Method: entry.Method,
				Count:  1,
			}
		}
	}

	if totalRequests == 0 {
		return nil, fmt.Errorf("未找到IP %s 的访问记录", targetIP)
	}

	// 分离GET和POST请求
	getURLs := make([]*models.URLInfo, 0)
	postURLs := make([]*models.URLInfo, 0)

	for _, info := range urlMap {
		if info.Method == "GET" {
			getURLs = append(getURLs, info)
		} else if info.Method == "POST" {
			postURLs = append(postURLs, info)
		}
	}

	// 排序
	sort.Slice(getURLs, func(i, j int) bool {
		return getURLs[i].Count > getURLs[j].Count
	})
	sort.Slice(postURLs, func(i, j int) bool {
		return postURLs[i].Count > postURLs[j].Count
	})

	// 生成TOP10
	getTop10 := as.generateURLTop10ForIP(getURLs, totalRequests)
	postTop10 := as.generateURLTop10ForIP(postURLs, totalRequests)

	// 判断IP类型
	ipType := "外网IP"
	if as.isInternalIP(targetIP) {
		ipType = "内网IP"
	}

	result := &models.SpecificIPResult{
		IP:            targetIP,
		IPType:        ipType,
		TotalRequests: totalRequests,
		FirstAccess:   firstAccess.Format("2006-01-02 15:04:05"),
		LastAccess:    lastAccess.Format("2006-01-02 15:04:05"),
		URLAnalysis: models.SpecificIPURLAnalysis{
			GetTop10:  getTop10,
			PostTop10: postTop10,
		},
	}

	return result, nil
}

// generateURLTop10ForIP 为特定IP生成URL TOP10
func (as *AnalysisService) generateURLTop10ForIP(urls []*models.URLInfo, totalRequests int) []models.URLStat {
	result := make([]models.URLStat, 0) // 确保初始化为空切片而不是nil
	limit := 10
	if len(urls) < limit {
		limit = len(urls)
	}

	for i := 0; i < limit; i++ {
		info := urls[i]
		percentage := float64(info.Count) / float64(totalRequests) * 100
		result = append(result, models.URLStat{
			Rank:       i + 1,
			URL:        info.URL,
			Count:      info.Count,
			Percentage: percentage,
		})
	}

	return result
}
