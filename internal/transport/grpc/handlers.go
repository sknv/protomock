package grpc

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func registerPackage(server *grpc.Server, pacakge Package) {
	for _, file := range pacakge.Files {
		registerFile(server, file)
	}
}

func registerFile(server *grpc.Server, file File) {
	for i := 0; i < file.ProtoFile.Services().Len(); i++ {
		service := file.ProtoFile.Services().Get(i)
		registerService(server, service)
	}
}

func registerService(server *grpc.Server, service protoreflect.ServiceDescriptor) {
	// Register each method in the service
	for i := 0; i < service.Methods().Len(); i++ {
		method := service.Methods().Get(i)
		methodName := string(method.Name())
		inputType := method.Input()
		outputType := method.Output()

		// Create a handler for the method.
		handler := func(ctx context.Context, req any) (any, error) {
			// Log the request
			log.Printf("Received request for method: %s", methodName)

			// Convert the request to a map[string]any.
			requestMap, err := dynamicMessageToMap(req.(*dynamicpb.Message))
			if err != nil {
				return nil, err
			}
			log.Printf("Request as map: %v", requestMap)

			// Create a map[string]any to simulate the response data
			responseData := map[string]any{
				"message": requestMap["name"],
				"details": map[string]any{
					"code":   200,
					"status": "OK",
				},
			}

			// Create a dynamic response message.
			resp, err := mapToDynamicMessage(outputType, responseData)
			if err != nil {
				return nil, err
			}

			return resp, nil
		}

		// Register the handler with the server
		server.RegisterService(&grpc.ServiceDesc{
			ServiceName: string(service.FullName()),
			HandlerType: (*any)(nil),
			Methods: []grpc.MethodDesc{
				{
					MethodName: methodName,
					Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
						// Decode the request.
						req := dynamicpb.NewMessage(inputType)
						if err := dec(req); err != nil {
							return nil, err
						}

						// Call the handler.
						return handler(ctx, req)
					},
				},
			},
		}, nil)
	}
}

// dynamicMessageToMap converts a dynamic message to a map[string]any.
func dynamicMessageToMap(msg *dynamicpb.Message) (map[string]any, error) {
	result := make(map[string]any)
	descriptor := msg.Descriptor()

	// Iterate over all fields in the message.
	for i := 0; i < descriptor.Fields().Len(); i++ {
		field := descriptor.Fields().Get(i)
		fieldName := string(field.Name())

		// Handle map fields.
		if field.IsMap() {
			mapValue := msg.Get(field).Map()
			mapResult := make(map[string]any)
			mapValue.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
				keyStr, err := convertMapKey(key)
				if err != nil {
					log.Printf("Failed to convert map key: %v", err)
					return false
				}
				convertedValue, err := convertFieldValue(field.MapValue(), value)
				if err != nil {
					log.Printf("Failed to convert map value: %v", err)
					return false
				}
				mapResult[keyStr] = convertedValue
				return true
			})
			result[fieldName] = mapResult
			continue
		}

		// Handle repeated fields.
		if field.IsList() {
			list := msg.Get(field).List()
			values := make([]any, list.Len())
			for j := 0; j < list.Len(); j++ {
				value := list.Get(j)
				convertedValue, err := convertFieldValue(field, value)
				if err != nil {
					return nil, fmt.Errorf("failed to convert repeated field %s: %w", fieldName, err)
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
			return nil, fmt.Errorf("failed to convert field %s: %w", fieldName, err)
		}
		result[fieldName] = convertedValue
	}

	return result, nil
}

// convertFieldValue converts a protoreflect.Value to a Go type based on the field's kind.
func convertFieldValue(field protoreflect.FieldDescriptor, value protoreflect.Value) (any, error) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return value.Bool(), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return int32(value.Int()), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return value.Int(), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return uint32(value.Uint()), nil
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
		nestedMsg := value.Message().Interface().(*dynamicpb.Message)
		return dynamicMessageToMap(nestedMsg)
	default:
		return nil, fmt.Errorf("unsupported field type: %v", field.Kind())
	}
}

// convertMapKey converts a protoreflect.MapKey to a string.
func convertMapKey(key protoreflect.MapKey) (string, error) {
	switch key.Interface().(type) {
	case string:
		return key.String(), nil
	case int32, int64, uint32, uint64:
		return fmt.Sprintf("%d", key.Int()), nil
	case bool:
		return fmt.Sprintf("%t", key.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported map key type: %T", key.Interface())
	}
}

//
//
//

// mapToDynamicMessage converts a map[string]any to a *dynamicpb.Message.
func mapToDynamicMessage(descriptor protoreflect.MessageDescriptor, data map[string]any) (*dynamicpb.Message, error) {
	msg := dynamicpb.NewMessage(descriptor)

	for key, value := range data {
		field := descriptor.Fields().ByName(protoreflect.Name(key))
		if field == nil {
			return nil, fmt.Errorf("field %s not found in message descriptor", key)
		}

		if err := setFieldValue(msg, field, value); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %w", key, err)
		}
	}

	return msg, nil
}

// setFieldValue sets the value of a field in a dynamic message.
func setFieldValue(msg *dynamicpb.Message, field protoreflect.FieldDescriptor, value any) error {
	switch {
	case field.IsMap():
		// Handle map fields
		mapValue, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map for field %s, got %T", field.Name(), value)
		}
		mapMsg := msg.NewField(field).Map()
		for k, v := range mapValue {
			keyValue := protoreflect.ValueOfString(k) // Assuming string keys
			fieldValue, err := convertToFieldValue(field.MapValue(), v)
			if err != nil {
				return fmt.Errorf("failed to convert map value for field %s: %w", field.Name(), err)
			}
			mapMsg.Set(keyValue.MapKey(), fieldValue)
		}
		msg.Set(field, protoreflect.ValueOfMap(mapMsg))

	case field.IsList():
		// Handle repeated fields.
		listValue, ok := value.([]any)
		if !ok {
			return fmt.Errorf("expected slice for field %s, got %T", field.Name(), value)
		}
		listMsg := msg.NewField(field).List()
		for _, v := range listValue {
			fieldValue, err := convertToFieldValue(field, v)
			if err != nil {
				return fmt.Errorf("failed to convert list value for field %s: %w", field.Name(), err)
			}
			listMsg.Append(fieldValue)
		}
		msg.Set(field, protoreflect.ValueOfList(listMsg))

	default:
		// Handle non-repeated, non-map fields.
		fieldValue, err := convertToFieldValue(field, value)
		if err != nil {
			return fmt.Errorf("failed to convert value for field %s: %w", field.Name(), err)
		}
		msg.Set(field, fieldValue)
	}

	return nil
}

// convertToFieldValue converts a Go value to a protoreflect.Value based on the field's kind.
func convertToFieldValue(field protoreflect.FieldDescriptor, value any) (protoreflect.Value, error) {
	switch field.Kind() {
	case protoreflect.BoolKind:
		v, ok := value.(bool)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected bool, got %T", value)
		}
		return protoreflect.ValueOfBool(v), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		switch v := value.(type) {
		case int32:
			return protoreflect.ValueOfInt32(v), nil
		case int:
			return protoreflect.ValueOfInt32(int32(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected int32, got %T", value)
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		switch v := value.(type) {
		case int64:
			return protoreflect.ValueOfInt64(v), nil
		case int:
			return protoreflect.ValueOfInt64(int64(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected int64, got %T", value)
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		switch v := value.(type) {
		case uint32:
			return protoreflect.ValueOfUint32(v), nil
		case int:
			return protoreflect.ValueOfUint32(uint32(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected uint32, got %T", value)
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		switch v := value.(type) {
		case uint64:
			return protoreflect.ValueOfUint64(v), nil
		case int:
			return protoreflect.ValueOfUint64(uint64(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected uint64, got %T", value)
		}
	case protoreflect.FloatKind:
		switch v := value.(type) {
		case float32:
			return protoreflect.ValueOfFloat32(v), nil
		case float64:
			return protoreflect.ValueOfFloat32(float32(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected float32, got %T", value)
		}
	case protoreflect.DoubleKind:
		switch v := value.(type) {
		case float64:
			return protoreflect.ValueOfFloat64(v), nil
		case float32:
			return protoreflect.ValueOfFloat64(float64(v)), nil
		default:
			return protoreflect.Value{}, fmt.Errorf("expected float64, got %T", value)
		}
	case protoreflect.StringKind:
		v, ok := value.(string)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected string, got %T", value)
		}
		return protoreflect.ValueOfString(v), nil
	case protoreflect.BytesKind:
		v, ok := value.([]byte)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected []byte, got %T", value)
		}
		return protoreflect.ValueOfBytes(v), nil
	case protoreflect.EnumKind:
		v, ok := value.(int32)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected int32 for enum, got %T", value)
		}
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)), nil
	case protoreflect.MessageKind:
		nestedMap, ok := value.(map[string]any)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("expected map for nested message, got %T", value)
		}
		nestedMsg, err := mapToDynamicMessage(field.Message(), nestedMap)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("failed to convert nested message: %w", err)
		}
		return protoreflect.ValueOfMessage(nestedMsg), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported field type: %v", field.Kind())
	}
}
