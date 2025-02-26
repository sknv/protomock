package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/dynamicpb"
)

type Handlers struct {
	packages Packages
}

func NewHandlers(packages Packages) *Handlers {
	return &Handlers{
		packages: packages,
	}
}

func (h *Handlers) Route(server *grpc.Server) {
	for _, pkg := range h.packages {
		registerPackage(server, pkg)
	}
}

func registerPackage(server *grpc.Server, pkg Package) {
	for _, file := range pkg.Files {
		registerFile(server, file)
	}
}

func registerFile(server *grpc.Server, file File) {
	for _, svc := range file.Services {
		registerService(server, svc)
	}
}

func registerService(server *grpc.Server, service Service) {
	// Register each method in the service.
	grpcMethods := make([]grpc.MethodDesc, 0, len(service.Mocks))

	for _, mock := range service.Mocks {
		method := mock.ProtoMethod
		methodName := string(method.Name())
		inputType := method.Input()
		outputType := method.Output()

		// Create a handler for the method.
		handler := func(ctx context.Context, req any) (any, error) {
			// Convert the request to a map[string]any.
			reqMessage, _ := req.(*dynamicpb.Message)

			request, err := NewMockRequestFrom(ctx, reqMessage)
			if err != nil {
				return nil, fmt.Errorf("decode request: %w", err)
			}

			response, err := mock.Eval(ctx, request)
			if err != nil {
				return nil, fmt.Errorf("evaluate mock: %w", err)
			}

			// Create a dynamic response message.
			reply, err := response.GRPC(outputType)
			if err != nil {
				return nil, fmt.Errorf("encode response: %w", err)
			}

			return reply, nil
		}

		grpcMethods = append(grpcMethods, grpc.MethodDesc{
			MethodName: methodName,
			Handler: func(
				srv any,
				ctx context.Context,
				decode func(any) error,
				intercept grpc.UnaryServerInterceptor,
			) (any, error) {
				// Decode the request.
				req := dynamicpb.NewMessage(inputType)
				if err := decode(req); err != nil {
					return nil, err
				}

				// Call the interceptor.
				if intercept != nil {
					return intercept(ctx, req, &grpc.UnaryServerInfo{
						Server:     srv,
						FullMethod: string(method.FullName()),
					}, func(ctx context.Context, req any) (any, error) {
						// Call the handler.
						return handler(ctx, req)
					})
				}

				// If no interceptor is provided, call the handler directly.
				return handler(ctx, req)
			},
		})
	}

	// Register the handler with the server.
	server.RegisterService(&grpc.ServiceDesc{ //nolint:exhaustruct // only required fields
		ServiceName: string(service.ProtoService.FullName()),
		HandlerType: (*any)(nil),
		Methods:     grpcMethods,
	}, nil)
}
