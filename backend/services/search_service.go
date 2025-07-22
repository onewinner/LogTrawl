package services

import (
	"regexp"
	"strings"

	"LogTrawl/backend/models"
)

type SearchService struct {
	fileService *FileService
}

func NewSearchService(fileService *FileService) *SearchService {
	return &SearchService{
		fileService: fileService,
	}
}

// SearchInFile searches for a pattern in the log file
func (ss *SearchService) SearchInFile(filePath, pattern string, caseSensitive, isRegex bool) ([]models.SearchResult, error) {
	content, err := ss.fileService.ReadLogFile(filePath)
	if err != nil {
		return nil, err
	}

	if isRegex {
		return ss.searchWithRegex(content.Lines, pattern, caseSensitive)
	}

	return ss.searchWithString(content.Lines, pattern, caseSensitive)
}

func (ss *SearchService) searchWithString(lines []string, pattern string, caseSensitive bool) ([]models.SearchResult, error) {
	var results []models.SearchResult
	searchPattern := pattern
	if !caseSensitive {
		searchPattern = strings.ToLower(pattern)
	}

	for i, line := range lines {
		searchLine := line
		if !caseSensitive {
			searchLine = strings.ToLower(line)
		}

		if strings.Contains(searchLine, searchPattern) {
			// 找到匹配位置
			var matches []int
			start := 0
			for {
				index := strings.Index(searchLine[start:], searchPattern)
				if index == -1 {
					break
				}
				matches = append(matches, start+index)
				start += index + len(searchPattern)
			}

			results = append(results, models.SearchResult{
				LineNumber: i + 1,
				Content:    line,
				Matches:    matches,
			})
		}
	}

	return results, nil
}

func (ss *SearchService) searchWithRegex(lines []string, pattern string, caseSensitive bool) ([]models.SearchResult, error) {
	flags := ""
	if !caseSensitive {
		flags = "(?i)"
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		return []models.SearchResult{}, err
	}

	var results []models.SearchResult
	for i, line := range lines {
		matches := regex.FindAllStringIndex(line, -1)
		if len(matches) > 0 {
			var matchPositions []int
			for _, match := range matches {
				matchPositions = append(matchPositions, match[0])
			}

			results = append(results, models.SearchResult{
				LineNumber: i + 1,
				Content:    line,
				Matches:    matchPositions,
			})
		}
	}

	return results, nil
}
