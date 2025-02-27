package dynamic

import (
	"fmt"

	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func MapToMessage(descriptor protoreflect.MessageDescriptor, data map[string]any) (*dynamicpb.Message, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encode data to json: %w", err)
	}

	msg := dynamicpb.NewMessage(descriptor)
	if err = protojson.Unmarshal(jsonData, msg); err != nil {
		return nil, fmt.Errorf("decode proto message from json: %w", err)
	}

	return msg, nil
}
