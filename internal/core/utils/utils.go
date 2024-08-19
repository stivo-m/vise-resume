package utils

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stivo-m/vise-resume/internal/core/dto"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func FormatApiResponse(message string, data interface{}) dto.ApiResponse[any] {
	var errorMap dto.ApiResponse[any]
	errorMap.Message = message
	errorMap.Data = data

	return errorMap
}

func GetValidationMessage(err validator.FieldError, field string) string {
	switch err.Tag() {
	case "required":
		return "The '" + field + "' field is required"
	case "min":
		return "The '" + field + "' field must be at least " + err.Param() + " characters long"
	case "oneof":
		return "The '" + field + "' field should be one of " + err.Param()
	default:
		return "The '" + field + "' field is invalid"
	}
}

func GetJSONFieldName(dto interface{}, structField string) string {
	r := reflect.TypeOf(dto).Elem()
	field, _ := r.FieldByName(structField)
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return strings.ToLower(structField)
	}
	return strings.Split(jsonTag, ",")[0]
}

func ExtractUuidFromContext(ctx context.Context, key string) (*uuid.UUID, error) {
	id, ok := ctx.Value(key).(string)

	if !ok {
		return nil, errors.New("given key is not set")
	}

	userId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	return &userId, nil
}
