
// We need a utility that takes the raw text from GitHub and splits it into logical files.
package service

import (
	"path/filepath"
	"strings"
)

// FileChange represents a single file modified in a PR
type FileChange struct {
	Path     string // e.g., "backend/main.go"
	Language string // e.g., "go", "python"
	Content  string // The actual code changes (diff patch)
	IsSafe   bool   // false if it's a lockfile, image, etc.
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

		// 1. Filter Junk
		if isIgnoredFile(path) {
			continue
		}

		// 2. Detect Language
		lang := detectLanguage(path)

		// 3. HANDLE LARGE FILES (Roadmap Checkbox ✅)
		content := "diff --git " + rawFile
		if len(content) > MaxFileSize {
			// Option A: Truncate it (Keep top 20KB, discard rest)
			content = content[:MaxFileSize] + "\n... [TRUNCATED DUE TO SIZE] ..."
			
			// Option B: Skip it entirely (Uncomment below if you prefer skipping)
			// continue 
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
		// Line 0 looks like: "a/filename.ext b/filename.ext"
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			// Return the second part (b/filename.ext), removing "b/" prefix
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
		"package-lock.json", "yarn.lock", "go.sum", // Lockfiles
		".png", ".jpg", ".svg", ".ico",             // Images
		".gitignore", ".env",                       // Configs
		"dist/", "build/", "node_modules/",         // Generated folders
	}

	for _, ignore := range ignored {
		if strings.Contains(path, ignore) || strings.HasSuffix(path, ignore) {
			return true
		}
	}
	return false
}