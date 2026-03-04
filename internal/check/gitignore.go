package check

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type GitignoreCheck struct {
	Dir string
}

func (c *GitignoreCheck) Name() string {
	return ".gitignore covers sensitive files"
}

func (c *GitignoreCheck) Run(_ context.Context) Result {
	gitignorePath := c.Dir + "/.gitignore"
	
	// Read .gitignore content
	file, err := os.Open(gitignorePath)
	if err != nil {
		return Result{
			Name:    c.Name(),
			Status:  StatusFail,
			Message: ".gitignore not found",
			Fix:     "Create a .gitignore file with common patterns",
		}
	}
	defer file.Close()

	// Parse .gitignore and collect all patterns
	patterns := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns[line] = true
	}

	// Define required patterns
	requiredPatterns := map[string]string{
		".env":          "Environment files with secrets",
		"*.log":         "Log files",
		"__pycache__":   "Python cache",
		"node_modules":  "Node.js dependencies",
	}

	// Check which patterns are missing
	var missing []string
	for pattern, description := range requiredPatterns {
		// Check if pattern or a variant exists
		found := false
		for existingPattern := range patterns {
			// Match exact pattern or with trailing slash (for directories)
			if existingPattern == pattern || 
			   existingPattern == pattern+"/" ||
			   existingPattern == "**/"+pattern ||
			   existingPattern == "**/"+pattern+"/" {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, fmt.Sprintf("%s (%s)", pattern, description))
		}
	}

	if len(missing) > 0 {
		return Result{
			Name:    c.Name(),
			Status:  StatusWarn,
			Message: fmt.Sprintf("Missing %d sensitive file patterns", len(missing)),
			Fix:     fmt.Sprintf("Add these patterns to .gitignore:\n  %s", strings.Join(missing, "\n  ")),
		}
	}

	return Result{
		Name:    c.Name(),
		Status:  StatusPass,
		Message: "Common sensitive files are gitignored",
	}
}
