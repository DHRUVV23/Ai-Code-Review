
// We need a utility that takes the raw text from GitHub and splits it into logical files.
package service

import (
	"path/filepath"
	"strings"
)

// FileChange represents a single file modified in a PR
type FileChange struct {
	Path     string 
	Language string 
	Content  string 
	IsSafe   bool  
}

const MaxFileSize = 20000

type DiffParser struct{}

func NewDiffParser() *DiffParser {
	return &DiffParser{}
}

// Parse splits a raw diff string into a list of FileChange objects
func (p *DiffParser) Parse(rawDiff string) []FileChange {
	var files []FileChange
	
	rawFiles := strings.Split(rawDiff, "diff --git ")

	for _, rawFile := range rawFiles {
		if strings.TrimSpace(rawFile) == "" {
			continue
		}

		path := extractFilePath(rawFile)
		if path == "" {
			continue
		}

		//  Filter Junk
		if isIgnoredFile(path) {
			continue
		}

		//  Detect Language
		lang := detectLanguage(path)

		//  HANDLE LARGE FILES
		content := "diff --git " + rawFile
		if len(content) > MaxFileSize {
		
			content = content[:MaxFileSize] + "\n... [TRUNCATED DUE TO SIZE] ..."
			
		}

		files = append(files, FileChange{
			Path:     path,
			Language: lang,
			Content:  content,
			IsSafe:   true,
		})
	}

	return files
}
// extractFilePath finds "a/backend/main.go b/backend/main.go" and returns "backend/main.go"
func extractFilePath(rawChunk string) string {
	lines := strings.Split(rawChunk, "\n")
	if len(lines) > 0 {
		
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
		
			return strings.TrimPrefix(parts[1], "b/")
		}
	}
	return ""
}

// detectLanguage guesses language based on extension
func detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "Go"
	case ".js", ".jsx", ".ts", ".tsx":
		return "JavaScript/TypeScript"
	case ".py":
		return "Python"
	case ".java":
		return "Java"
	case ".html", ".css":
		return "Web"
	case ".md":
		return "Markdown"
	default:
		return "Unknown"
	}
}

// isIgnoredFile returns true if we should skip sending this to AI
func isIgnoredFile(path string) bool {
	ignored := []string{
		"package-lock.json", "yarn.lock", "go.sum", 
		".png", ".jpg", ".svg", ".ico",            
		".gitignore", ".env",                       
		"dist/", "build/", "node_modules/",         
	}

	for _, ignore := range ignored {
		if strings.Contains(path, ignore) || strings.HasSuffix(path, ignore) {
			return true
		}
	}
	return false
}