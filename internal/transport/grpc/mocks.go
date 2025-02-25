package grpc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
)

const (
	_mockFileExtension  = ".js"
	_protoFileExtension = ".proto"
)

// type Mock struct {
// 	Method string
// 	Script string
// }

// type Mocks []Mock

// type Service struct {
// 	Mocks Mocks
// }

// type Services []Service

type File struct {
	ProtoFile linker.File
	// Services  Service
}

type Files []File

type Package struct {
	Files Files
}

type Packages []Package

// ----------------------------------------------------------------------------

// BuildPackages traverses the directory and populate Packages.
func BuildPackages(ctx context.Context, mocksDir string) (Packages, error) {
	protoFiles := make(map[string]Files) // Map of package to its files.
	compiler := &protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{},
	}

	// Walk through the directory
	err := filepath.Walk(mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("traverse path: %w", err)
		}

		// Skip directories.
		if info.IsDir() {
			return nil
		}

		// Process only files with the proper extensions.
		ext := filepath.Ext(path)

		switch ext {
		case _mockFileExtension:
			return nil
		case _protoFileExtension:
			protoFile, err := buildProtoFile(ctx, compiler, path)
			if err != nil {
				return fmt.Errorf("build proto file: %w", err)
			}

			packageName := string(protoFile.Package())
			protoFiles[packageName] = append(protoFiles[packageName], File{
				ProtoFile: protoFile,
			})

			return nil
		default:
			return nil
		}
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	packages := make(Packages, 0, len(protoFiles))
	for _, files := range protoFiles {
		packages = append(packages, Package{
			Files: files,
		})
	}

	return packages, nil
}

func buildProtoFile(ctx context.Context, compiler *protocompile.Compiler, path string) (linker.File, error) {
	files, err := compiler.Compile(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("compile proto file: %w", err)
	}

	return files[0], nil // We only process one file at a time.
}
