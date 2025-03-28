package http

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sknv/protomock/pkg/js"
	xstrings "github.com/sknv/protomock/pkg/strings"
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

func (m Mock) Eval(ctx context.Context, request MockRequest) (MockResponse, error) {
	vm := js.NewRuntime()

	console := js.NewConsole(ctx)
	if err := vm.Set("console", console); err != nil {
		return MockResponse{}, fmt.Errorf("set console in runtime: %w", err)
	}

	if err := vm.Set("request", request); err != nil {
		return MockResponse{}, fmt.Errorf("set request in runtime: %w", err)
	}

	eval, err := vm.RunString(m.Script)
	if err != nil {
		return MockResponse{}, fmt.Errorf("eval script: %w", err)
	}

	var response MockResponse
	if err = vm.ExportTo(eval, &response); err != nil {
		return MockResponse{}, fmt.Errorf("export response from js: %w", err)
	}

	return response, nil
}

// ----------------------------------------------------------------------------

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

		// Build a mock from file.
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		httpMethod := strings.TrimSuffix(info.Name(), _mockFileExtension)
		httpPath := filepath.ToSlash( // Transform OS-specific separators to slashes.
			strings.TrimPrefix(filepath.Dir(path), filepath.Clean(mocksDir)), // Trim original mocks dir.
		)
		httpPath = strings.ReplaceAll( // Replace wildcards for router.
			httpPath,
			wildcardPatternToFind,
			wildcardPaternToReplace,
		)

		mock := Mock{
			Method: httpMethod,
			Path:   httpPath,
			Script: xstrings.ByteSliceToString(content),
		}

		mocks = append(mocks, mock)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	return mocks, nil
}
