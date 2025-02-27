package grpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/sknv/protomock/pkg/protobuf/dynamic"
)

type MockResponseBody map[string]any

type MockResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MockResponse struct {
	Body  MockResponseBody   `json:"body"`
	Error *MockResponseError `json:"error"`
}

func (r MockResponse) GRPC(response protoreflect.MessageDescriptor) (*dynamicpb.Message, error) {
	if r.Error != nil {
		// Create a gRPC status with error details if needed.
		sts := status.New(codes.Code(r.Error.Code), r.Error.Message) //nolint:gosec // determined range

		return nil, sts.Err() //nolint:wrapcheck // plain gRPC error
	}

	// If there is no error, return a response.
	message, err := dynamic.MapToMessage(response, r.Body)
	if err != nil {
		return nil, fmt.Errorf("encode proto body: %w", err)
	}

	return message, nil
}
