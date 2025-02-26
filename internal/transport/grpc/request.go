package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/sknv/protomock/pkg/protobuf/dynamic"
)

type (
	MockRequestMetadata map[string]string
	MockRequestBody     map[string]any
)

type MockRequest struct {
	Metadata MockRequestMetadata `json:"metadata"`
	Body     MockRequestBody     `json:"body"`
}

func NewMockRequestFrom(ctx context.Context, r *dynamicpb.Message) (MockRequest, error) {
	body, err := dynamic.MessageToMap(r)
	if err != nil {
		return MockRequest{}, fmt.Errorf("decode proto body: %w", err)
	}

	md, _ := metadata.FromIncomingContext(ctx)
	meta := make(MockRequestMetadata, len(md))

	for key := range md {
		var firstVal string

		values := md.Get(key)
		if len(values) > 0 {
			firstVal = values[0]
		}

		meta[key] = firstVal
	}

	return MockRequest{
		Metadata: meta,
		Body:     body,
	}, nil
}
