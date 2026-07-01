package middleware

import (
	"context"
	"fmt"
	"reflect"
	"time"

	loggerPkg "git.itechpsp.com/e46/box/platform.git/modules/logger"
	server "git.itechpsp.com/e46/box/platform.git/pkg/micro/server"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func LoggingWrapper(logger loggerPkg.Logger) server.HandlerWrapper {
	return func(handler server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			startTime := time.Now()

			logRequest(ctx, logger, req.Method(), req.Body())

			err := handler(ctx, req, rsp)

			duration := time.Since(startTime)

			logResponse(ctx, logger, req.Method(), rsp, err, duration)

			return err
		}
	}
}

func logRequest(ctx context.Context, logger loggerPkg.Logger, method string, req interface{}) {
	if reqJSON, err := MarshalWithReadableConversion(req); err == nil {
		logger.For(ctx).Info(fmt.Sprintf("gRPC Request [%s]:\n%s", method, string(reqJSON)))
	} else {
		logger.For(ctx).Info(fmt.Sprintf("gRPC Request [%s]: %+v", method, req))
	}
}

func logResponse(ctx context.Context, logger loggerPkg.Logger, method string, resp interface{}, err error, duration time.Duration) {
	if err != nil {
		logger.For(ctx).Error(fmt.Sprintf("gRPC Response [%s] ERROR (duration: %s): %v", method, duration.String(), err))
		return
	}

	if respJSON, err := MarshalWithReadableConversion(resp); err == nil {
		logger.For(ctx).Info(fmt.Sprintf("gRPC Response [%s] (duration: %s):\n%s", method, duration.String(), string(respJSON)))
	} else {
		logger.For(ctx).Info(fmt.Sprintf("gRPC Response [%s] (duration: %s): %+v", method, duration.String(), resp))
	}
}

func MarshalWithReadableConversion(v interface{}) ([]byte, error) {
	converted := ConvertToReadable(v)

	return json.MarshalIndent(converted, "", "  ")
}

func ConvertToReadable(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		result := make(map[string]interface{})
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			if !field.CanInterface() {
				continue
			}

			fieldName := fieldType.Name
			fieldValue := field.Interface()

			if converted := convertField(fieldValue); converted != nil {
				result[fieldName] = converted
			} else {
				result[fieldName] = ConvertToReadable(fieldValue)
			}
		}
		return result

	case reflect.Slice:
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = ConvertToReadable(val.Index(i).Interface())
		}
		return result

	case reflect.Map:
		result := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			result[keyStr] = ConvertToReadable(val.MapIndex(key).Interface())
		}
		return result

	default:
		return v
	}
}

func convertField(v interface{}) interface{} {
	switch enum := v.(type) {
	case *timestamppb.Timestamp:
		if enum == nil {
			return nil
		}
		return enum.AsTime().Format(time.RFC3339)
	case timestamppb.Timestamp:
		return enum.AsTime().Format(time.RFC3339)
	default:
		return nil
	}
}
