package grpc

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/sknv/protomock/pkg/protobuf/dynamic"
)

type (
	MockResponseBody map[string]any
)

type MockResponse struct {
	Body MockResponseBody `json:"body"`
}

func (r MockResponse) GRPC(response protoreflect.MessageDescriptor) (*dynamicpb.Message, error) {
	message, err := dynamic.MapToMessage(response, r.Body)
	if err != nil {
		return nil, fmt.Errorf("encode proto body: %w", err)
	}

	return message, nil
}
