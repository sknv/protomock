package dynamic

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// MapToMessage converts a map[string]any to a *dynamicpb.Message.
func MapToMessage(descriptor protoreflect.MessageDescriptor, data map[string]any) (*dynamicpb.Message, error) {
	msg := dynamicpb.NewMessage(descriptor)

	for key, value := range data {
		field := descriptor.Fields().ByName(protoreflect.Name(key))
		if field == nil {
			return nil, fmt.Errorf("field %s not found in message descriptor", key) //nolint:err113
		}

		if err := setFieldValue(msg, field, value); err != nil {
			return nil, fmt.Errorf("set field %s: %w", key, err)
		}
	}

	return msg, nil
}

// setFieldValue sets the value of a field in a dynamic message.
//
//nolint:cyclop // helper function
func setFieldValue(msg *dynamicpb.Message, field protoreflect.FieldDescriptor, value any) error {
	switch {
	case field.IsMap():
		// Handle map fields.
		mapValue, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map for field %s, got %T", field.Name(), value) //nolint:err113
		}

		mapMsg := msg.NewField(field).Map()

		for k, v := range mapValue {
			keyValue := protoreflect.ValueOfString(k) // Assuming string keys.

			fieldValue, err := convertToFieldValue(field.MapValue(), v)
			if err != nil {
				return fmt.Errorf("convert map value for field %s: %w", field.Name(), err)
			}

			mapMsg.Set(keyValue.MapKey(), fieldValue)
		}

		msg.Set(field, protoreflect.ValueOfMap(mapMsg))
	case field.IsList():
		// Handle repeated fields.
		listValue, ok := value.([]any)
		if !ok {
			return fmt.Errorf("expected slice for field %s, got %T", field.Name(), value) //nolint:err113
		}

		listMsg := msg.NewField(field).List()

		for _, v := range listValue {
			fieldValue, err := convertToFieldValue(field, v)
			if err != nil {
				return fmt.Errorf("convert list value for field %s: %w", field.Name(), err)
			}

			listMsg.Append(fieldValue)
		}

		msg.Set(field, protoreflect.ValueOfList(listMsg))
	default:
		// Handle non-repeated, non-map fields.
		fieldValue, err := convertToFieldValue(field, value)
		if err != nil {
			return fmt.Errorf("convert value for field %s: %w", field.Name(), err)
		}

		msg.Set(field, fieldValue)
	}

	return nil
}

// convertToFieldValue converts a Go value to a protoreflect.Value based on the field's kind.
//
//nolint:cyclop,funlen,gocyclo // helper function
func convertToFieldValue(field protoreflect.FieldDescriptor, value any) (protoreflect.Value, error) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		v, ok := value.(bool)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected bool, got %T", value) //nolint:err113
		}

		return protoreflect.ValueOfBool(v), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		switch v := value.(type) {
		case int32:
			return protoreflect.ValueOfInt32(v), nil
		case int: // Default integer type.
			return protoreflect.ValueOfInt32(int32(v)), nil //nolint:gosec // possible overflow
		case int64:
			return protoreflect.ValueOfInt32(int32(v)), nil //nolint:gosec // possible overflow
		default:
			return protoreflect.Value{}, fmt.Errorf("expected int32, got %T", value) //nolint:err113
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		switch v := value.(type) {
		case int64:
			return protoreflect.ValueOfInt64(v), nil
		case int: // Default integer type.
			return protoreflect.ValueOfInt64(int64(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected int64, got %T", value) //nolint:err113
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		switch v := value.(type) {
		case uint32:
			return protoreflect.ValueOfUint32(v), nil
		case int: // Default integer type.
			return protoreflect.ValueOfUint32(uint32(v)), nil //nolint:gosec // possible overflow
		case int64:
			return protoreflect.ValueOfUint32(uint32(v)), nil //nolint:gosec // possible overflow
		default:
			return protoreflect.Value{}, fmt.Errorf("expected uint32, got %T", value) //nolint:err113
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		switch v := value.(type) {
		case uint64:
			return protoreflect.ValueOfUint64(v), nil
		case int: // Default integer type.
			return protoreflect.ValueOfUint64(uint64(v)), nil //nolint:gosec // possible overflow
		case int64:
			return protoreflect.ValueOfUint64(uint64(v)), nil //nolint:gosec // possible overflow
		default:
			return protoreflect.Value{}, fmt.Errorf("expected uint64, got %T", value) //nolint:err113
		}
	case protoreflect.FloatKind:
		switch v := value.(type) {
		case float32:
			return protoreflect.ValueOfFloat32(v), nil
		case float64: // Default float type.
			return protoreflect.ValueOfFloat32(float32(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected float32, got %T", value) //nolint:err113
		}
	case protoreflect.DoubleKind:
		v, ok := value.(float64)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected float64, got %T", value) //nolint:err113
		}

		return protoreflect.ValueOfFloat64(v), nil
	case protoreflect.StringKind:
		v, ok := value.(string)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected string, got %T", value) //nolint:err113
		}

		return protoreflect.ValueOfString(v), nil
	case protoreflect.BytesKind:
		v, ok := value.([]byte)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected []byte, got %T", value) //nolint:err113
		}

		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.EnumKind:
		v, ok := value.(int32)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected int32 for enum, got %T", value) //nolint:err113
		}

		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)), nil
	case protoreflect.MessageKind:
		nestedMap, ok := value.(map[string]any)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected map for nested message, got %T", value) //nolint:err113
		}

		nestedMsg, err := MapToMessage(field.Message(), nestedMap)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("convert nested message: %w", err)
		}

		return protoreflect.ValueOfMessage(nestedMsg), nil
	case protoreflect.GroupKind:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %v", field.Kind()) //nolint:err113
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %v", field.Kind()) //nolint:err113
	}
}
