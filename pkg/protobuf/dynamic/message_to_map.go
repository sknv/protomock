package dynamic

import (
	"fmt"

	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
)

func MessageToMap(msg *dynamicpb.Message) (map[string]any, error) {
	jsonData, err := protojson.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("encode proto message to json: %w", err)
	}

	var data map[string]any
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("decode data from json: %w", err)
	}

	return data, nil
}
