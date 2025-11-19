package grpc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sknv/protomock/pkg/js"
	xstrings "github.com/sknv/protomock/pkg/strings"
)

const (
	_mockFileExtension  = ".js"
	_protoFileExtension = ".proto"

	_protoIncludePath = "./include"
)

type Mock struct {
	ProtoMethod protoreflect.MethodDescriptor
	Script      string
}

type Mocks []Mock

type Service struct {
	ProtoService protoreflect.ServiceDescriptor
	Mocks        Mocks
}

type Services []Service

type File struct {
	ProtoFile linker.File
	Services  Services
}

type Files []File

type Package struct {
	Files Files
}

type Packages []Package

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

type mockID struct {
	Package string
	Service string
	Method  string
}

// BuildPackages traverses the directory and populate Packages.
//
//nolint:funlen // mostly basic operations
func BuildPackages(ctx context.Context, mocksDir string) (Packages, error) {
	var (
		mocks = make(map[mockID]Mock)

		protoFiles []linker.File
	)

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
			// Build a mock from file.
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}

			method := strings.TrimSuffix(info.Name(), _mockFileExtension)
			service := filepath.Base(filepath.Dir(path))
			pkg := filepath.ToSlash( // Transform OS-specific separators to slashes.
				strings.TrimPrefix(filepath.Dir(path), filepath.Clean(mocksDir)), // Trim original mocks dir.
			)
			pkg = strings.TrimSuffix(pkg, service) // Trim service name.
			pkg = strings.Trim(pkg, "/")           // Get rid of path slashes.

			mockID := mockID{
				Package: pkg,
				Service: service,
				Method:  method,
			}
			mock := Mock{
				ProtoMethod: nil, // Will be mapped later.
				Script:      xstrings.ByteSliceToString(content),
			}

			mocks[mockID] = mock

			return nil
		case _protoFileExtension:
			curDir := filepath.Dir(path)
			curFile := filepath.Base(path)

			// Parse proto definition.
			protoFile, err := buildProtoFile(ctx, []string{_protoIncludePath, curDir}, curFile)
			if err != nil {
				return fmt.Errorf("build proto file: %w", err)
			}

			protoFiles = append(protoFiles, protoFile)

			return nil
		default:
			return nil
		}
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	return mapProtoFilesToMocks(protoFiles, mocks), nil
}

//nolint:ireturn,nolintlint // contract
func buildProtoFile(
	ctx context.Context,
	importPaths []string,
	file string,
) (linker.File, error) {
	//nolint:exhaustruct // only required field
	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: importPaths,
		},
	}

	files, err := compiler.Compile(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("compile proto file: %w", err)
	}

	return files[0], nil // We only process one file at a time.
}

func mapProtoFilesToMocks(protoFiles linker.Files, mocks map[mockID]Mock) Packages {
	files := make(map[string]Files) // Map of package name to files.

	for _, protoFile := range protoFiles {
		packageName := string(protoFile.Package())
		file := mapProtoFileToMocks(protoFile, mocks)

		files[packageName] = append(files[packageName], file)
	}

	pkgs := make(Packages, 0, len(files))
	for _, files := range files {
		pkgs = append(pkgs, Package{
			Files: files,
		})
	}

	return pkgs
}

func mapProtoFileToMocks(protoFile linker.File, mocks map[mockID]Mock) File {
	services := make(Services, 0, protoFile.Services().Len())

	for i := range protoFile.Services().Len() {
		protoService := protoFile.Services().Get(i)
		service := mapProtoServiceToMocks(string(protoFile.Package()), protoService, mocks)

		services = append(services, service)
	}

	return File{
		ProtoFile: protoFile,
		Services:  services,
	}
}

func mapProtoServiceToMocks(
	packageName string,
	protoService protoreflect.ServiceDescriptor,
	mocks map[mockID]Mock,
) Service {
	var serviceMocks Mocks

	for i := range protoService.Methods().Len() {
		protoMethod := protoService.Methods().Get(i)
		mockID := mockID{
			Package: packageName,
			Service: string(protoService.Name()),
			Method:  string(protoMethod.Name()),
		}

		if svcMock, ok := mocks[mockID]; ok {
			svcMock.ProtoMethod = protoMethod
			serviceMocks = append(serviceMocks, svcMock)
		}
	}

	return Service{
		ProtoService: protoService,
		Mocks:        serviceMocks,
	}
}
