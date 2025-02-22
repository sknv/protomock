package http

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	_mockFileExtension = ".js"

	wildcardPatternToFind   = "__"
	wildcardPaternToReplace = ":"
)

type Mock struct {
	Method string
	Path   string
	Script string
}

type Mocks []Mock

// BuildMocks traverses the directory and populate Mocks.
func BuildMocks(mocksDir string) (Mocks, error) {
	var mocks Mocks

	// Walk through the directory
	err := filepath.Walk(mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traverse path: %w", err)
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Process only files with the proper extension.
		if filepath.Ext(path) != _mockFileExtension {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		// Build a mock from file.
		httpMethod := strings.TrimSuffix(info.Name(), _mockFileExtension)
		httpPath := strings.ReplaceAll( // Replace wildcards for router.
			filepath.ToSlash( // Transform OS-specific separators to slashes.
				strings.TrimPrefix(filepath.Dir(path), filepath.Clean(mocksDir)), // Trim original mocks dir.
			),
			wildcardPatternToFind,
			wildcardPaternToReplace,
		)

		mock := Mock{
			Method: httpMethod,
			Path:   httpPath,
			Script: string(content),
		}

		mocks = append(mocks, mock)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	return mocks, nil
}
