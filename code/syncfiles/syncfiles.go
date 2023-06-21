package syncfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	filePattern = regexp.MustCompile(`(?s)//// FILE~(?P<filepath>.*?) ////\n(?P<content>.*?)\n//// END FILE ////`)
	ignoreList  = []string{
		"messages.json",
		".git",
	}
)

func Update(parentFolder, update string) error {
	if !filepath.IsAbs(parentFolder) {
		return fmt.Errorf("parent folder path: %q is not absolute", parentFolder)
	}

	matches := filePattern.FindAllStringSubmatch(update, -1)
	for _, match := range matches {
		filePath := filepath.Join(parentFolder, match[1])
		content := match[2]

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directories for file: %q: %v", filePath, err)
		}

		if err := os.WriteFile(filePath, []byte(content), os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file: %q: %v", filePath, err)
		}
	}

	return nil
}

func Load(parentFolder string) (string, error) {
	if !filepath.IsAbs(parentFolder) {
		return "", fmt.Errorf("parent folder path: %q is not an absolute path", parentFolder)
	}

	var result strings.Builder
	err := filepath.Walk(parentFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access file: %q: %v", path, err)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		// Check ignore list
		for _, ignore := range ignoreList {
			if strings.Contains(path, ignore) {
				return nil
			}
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %q: %v", path, err)
		}

		relativePath, err := filepath.Rel(parentFolder, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for file: %q: %v", path, err)
		}

		result.WriteString(fmt.Sprintf("//// FILE~%s ////\n%s\n//// END FILE ////", relativePath, content))
		return nil
	})

	if err != nil {
		return "", err
	}

	return result.String(), nil
}
