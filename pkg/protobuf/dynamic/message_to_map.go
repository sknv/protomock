package dynamic

import (
	"fmt"
	"strconv"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// MessageToMap converts a dynamic message to a map[string]any.
//
//nolint:funlen // keep all in one place
func MessageToMap(msg *dynamicpb.Message) (map[string]any, error) {
	result := make(map[string]any)
	descriptor := msg.Descriptor()

	// Iterate over all fields in the message.
	for i := range descriptor.Fields().Len() {
		field := descriptor.Fields().Get(i)
		fieldName := string(field.Name())

		// Handle map fields.
		if field.IsMap() {
			var err error

			mapResult := make(map[string]any)
			mapValue := msg.Get(field).Map()

			mapValue.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
				keyStr, errConvert := convertMapKey(key)
				if errConvert != nil {
					err = fmt.Errorf("convert map field: %w", errConvert)

					return false
				}

				convertedValue, errConvert := convertFieldValue(field.MapValue(), value)
				if errConvert != nil {
					err = fmt.Errorf("convert map value: %w", errConvert)

					return false
				}

				mapResult[keyStr] = convertedValue

				return true
			})

			if err != nil {
				return nil, fmt.Errorf("convert map %s: %w", fieldName, err)
			}

			result[fieldName] = mapResult

			continue
		}

		// Handle repeated fields.
		if field.IsList() {
			list := msg.Get(field).List()
			values := make([]any, list.Len())

			for j := range list.Len() {
				value := list.Get(j)

				convertedValue, err := convertFieldValue(field, value)
				if err != nil {
					return nil, fmt.Errorf("convert repeated field %s: %w", fieldName, err)
				}

				values[j] = convertedValue
			}

			result[fieldName] = values

			continue
		}

		// Handle non-repeated, non-map fields.
		value := msg.Get(field)

		convertedValue, err := convertFieldValue(field, value)
		if err != nil {
			return nil, fmt.Errorf("convert field %s: %w", fieldName, err)
		}

		result[fieldName] = convertedValue
	}

	return result, nil
}

// convertFieldValue converts a protoreflect.Value to a Go type based on the field's kind.
//
//nolint:cyclop // helper function
func convertFieldValue(field protoreflect.FieldDescriptor, value protoreflect.Value) (any, error) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(value.Int()), nil //nolint:gosec // possible overflow
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(value.Uint()), nil //nolint:gosec // possible overflow
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return value.Uint(), nil
	case protoreflect.FloatKind:
		return float32(value.Float()), nil
	case protoreflect.DoubleKind:
		return value.Float(), nil
	case protoreflect.StringKind:
		return value.String(), nil
	case protoreflect.BytesKind:
		return value.Bytes(), nil
	case protoreflect.EnumKind:
		return value.Enum(), nil
	case protoreflect.MessageKind:
		nestedMsg, ok := value.Message().Interface().(*dynamicpb.Message)
		if !ok {
			return nil, fmt.Errorf("expected message type, got %T", value.Message().Interface()) //nolint:err113
		}

		return MessageToMap(nestedMsg)
	case protoreflect.GroupKind:
		return nil, fmt.Errorf("unsupported field type: %v", field.Kind()) //nolint:err113
	default:
		return nil, fmt.Errorf("unsupported field type: %v", field.Kind()) //nolint:err113
	}
}

// convertMapKey converts a protoreflect.MapKey to a string.
func convertMapKey(key protoreflect.MapKey) (string, error) {
	switch key.Interface().(type) {
	case string:
		return key.String(), nil
	case int32, int64, uint32, uint64:
		return strconv.FormatInt(key.Int(), 10), nil
	case bool:
		return strconv.FormatBool(key.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported proto map key type: %T", key.Interface()) //nolint:err113
	}
}
