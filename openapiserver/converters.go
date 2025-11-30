// Package openapiserver provides utilities for handling OpenAPI server requests and responses.
package openapiserver

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

// Error definitions for converter functions.
var (
	ErrFieldNotSettable         = errors.New("field is not settable")
	ErrOutMustBePointerToStruct = errors.New("out must be a non-nil pointer to struct")
	ErrOutMustBePointer         = errors.New("out must be a pointer to struct")
	ErrCannotSetField           = errors.New("cannot set field")
	ErrUnsupportedFieldType     = errors.New("unsupported field type")
	ErrUnsupportedSliceElemType = errors.New("unsupported slice element type")
)

func pathToStruct(r *http.Request, target any) error {
	v := reflect.ValueOf(target)

	v = v.Elem()
	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		structField := t.Field(i)

		// Get path tag value
		pathTag := structField.Tag.Get("path")
		if pathTag == "" {
			continue
		}

		// Get path parameter value
		paramValue := r.PathValue(pathTag)
		if paramValue == "" {
			continue
		}

		// Set field based on its type
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set: %w", structField.Name, ErrFieldNotSettable)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(paramValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intValue, err := strconv.ParseInt(paramValue, 10, 64); err == nil {
				field.SetInt(intValue)
			} else {
				return fmt.Errorf("invalid int value %s for field %s: %w", paramValue, structField.Name, err)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintValue, err := strconv.ParseUint(paramValue, 10, 64); err == nil {
				field.SetUint(uintValue)
			} else {
				return fmt.Errorf("invalid uint value %s for field %s: %w", paramValue, structField.Name, err)
			}

		case reflect.Bool:
			if boolValue, err := strconv.ParseBool(paramValue); err == nil {
				field.SetBool(boolValue)
			} else {
				return fmt.Errorf("invalid bool value %s for field %s: %w", paramValue, structField.Name, err)
			}

		case reflect.Float32, reflect.Float64:
			if floatValue, err := strconv.ParseFloat(paramValue, 64); err == nil {
				field.SetFloat(floatValue)
			} else {
				return fmt.Errorf("invalid float value %s for field %s: %w", paramValue, structField.Name, err)
			}

		default:
			return fmt.Errorf("unsupported field type %s for field %s: %w", field.Kind(), structField.Name, ErrUnsupportedFieldType)
		}
	}

	return nil
}

func mapFromStruct[T ~map[string][]string](in any, tag string) T {
	out := make(map[string][]string)
	v := reflect.ValueOf(in)
	t := v.Type()

	for i := range v.NumField() {
		field := t.Field(i)
		key := field.Tag.Get(tag)
		if key == "" {
			continue
		}
		if tag == "header" {
			key = http.CanonicalHeaderKey(key)
		}

		fieldValue := v.Field(i)
		var values []string

		// Check if the field is a slice
		if fieldValue.Kind() == reflect.Slice {
			for j := range fieldValue.Len() {
				elem := fieldValue.Index(j)
				values = append(values, fmt.Sprintf("%v", elem.Interface()))
			}
		} else {
			// Convert non-slice fields to string
			values = []string{fmt.Sprintf("%v", fieldValue.Interface())}
		}

		out[key] = values
	}
	return out
}

func mapToStruct[T ~map[string][]string](m T, tag string, out any) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return ErrOutMustBePointerToStruct
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrOutMustBePointer
	}

	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the header tag
		tag := fieldType.Tag.Get(tag)
		if tag == "" {
			continue
		}

		// Look up the value in the map
		values, exists := m[tag]
		if !exists || len(values) == 0 {
			continue
		}

		// Set the field value based on its type
		if err := setField(field, values); err != nil {
			return fmt.Errorf("field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setField sets a struct field value with enhanced type support
func setField(field reflect.Value, values []string) error {
	if !field.CanSet() {
		return ErrCannotSetField
	}

	fieldType := field.Type()

	switch fieldType.Kind() {
	case reflect.String:
		field.SetString(values[0])

	case reflect.Slice:
		return setSliceField(field, values)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(values) > 0 {
			intVal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int: %w", err)
			}
			field.SetInt(intVal)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(values) > 0 {
			uintVal, err := strconv.ParseUint(values[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse uint: %w", err)
			}
			field.SetUint(uintVal)
		}

	case reflect.Float32, reflect.Float64:
		if len(values) > 0 {
			floatVal, err := strconv.ParseFloat(values[0], 64)
			if err != nil {
				return fmt.Errorf("failed to parse float: %w", err)
			}
			field.SetFloat(floatVal)
		}

	case reflect.Bool:
		if len(values) > 0 {
			boolVal, err := strconv.ParseBool(values[0])
			if err != nil {
				return fmt.Errorf("failed to parse bool: %w", err)
			}
			field.SetBool(boolVal)
		}

	default:
		return fmt.Errorf("unsupported field type: %s: %w", fieldType.Kind(), ErrUnsupportedFieldType)
	}

	return nil
}

// setSliceField handles slice types
func setSliceField(field reflect.Value, values []string) error {
	elemType := field.Type().Elem()

	switch elemType.Kind() {
	case reflect.String:
		field.Set(reflect.ValueOf(values))

	case reflect.Int:
		intSlice := make([]int, len(values))
		for i, v := range values {
			intVal, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("failed to parse int in slice: %w", err)
			}
			intSlice[i] = intVal
		}
		field.Set(reflect.ValueOf(intSlice))

	case reflect.Int64:
		intSlice := make([]int64, len(values))
		for i, v := range values {
			intVal, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int64 in slice: %w", err)
			}
			intSlice[i] = intVal
		}
		field.Set(reflect.ValueOf(intSlice))

	case reflect.Float64:
		floatSlice := make([]float64, len(values))
		for i, v := range values {
			floatVal, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("failed to parse float64 in slice: %w", err)
			}
			floatSlice[i] = floatVal
		}
		field.Set(reflect.ValueOf(floatSlice))

	case reflect.Bool:
		boolSlice := make([]bool, len(values))
		for i, v := range values {
			boolVal, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("failed to parse bool in slice: %w", err)
			}
			boolSlice[i] = boolVal
		}
		field.Set(reflect.ValueOf(boolSlice))

	default:
		return fmt.Errorf("unsupported slice element type: %s: %w", elemType.Kind(), ErrUnsupportedSliceElemType)
	}

	return nil
}
