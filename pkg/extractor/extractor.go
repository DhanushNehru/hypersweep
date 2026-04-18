package extractor

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The URL regex captures basic domains and paths
// It limits characters to avoid trailing punctuation like ) or . in markdown.
var urlRegex = regexp.MustCompile(`https?://[^\s<>"'()]+(?:[a-zA-Z0-9/])`)

// Result represents an extracted URL and its file source
type Result struct {
	URL      string
	FilePath string
	LineNum  int
}

// Extractor handles finding all target files and extracting URLs
type Extractor struct {
	RootPath       string
	TargetExts     map[string]bool
	IgnoredDomains []string
}

// NewExtractor initializes an Extractor with default settings
func NewExtractor(rootPath string) *Extractor {
	return &Extractor{
		RootPath: rootPath,
		TargetExts: map[string]bool{
			".md":   true,
			".txt":  true,
			".html": true,
		},
		IgnoredDomains: []string{"localhost", "127.0.0.1"},
	}
}

// Extract starts the process and returns a slice of Results
func (e *Extractor) Extract() ([]Result, error) {
	var results []Result

	err := filepath.Walk(e.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip .git, node_modules, etc.
			if strings.HasPrefix(info.Name(), ".") || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !e.TargetExts[ext] {
			return nil
		}

		urls, err := e.extractFromFile(path)
		if err != nil {
			// Log warning maybe? For now just continue
			return nil
		}
		results = append(results, urls...)

		return nil
	})

	return results, err
}

func (e *Extractor) extractFromFile(path string) ([]Result, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileResults []Result
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		matches := urlRegex.FindAllString(line, -1)
		
		for _, match := range matches {
			if !e.isIgnored(match) {
				fileResults = append(fileResults, Result{
					URL:      match,
					FilePath: path,
					LineNum:  lineNum,
				})
			}
		}
		lineNum++
	}

	return fileResults, scanner.Err()
}

func (e *Extractor) isIgnored(url string) bool {
	for _, domain := range e.IgnoredDomains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}
