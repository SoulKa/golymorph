package objectpath

import (
	"fmt"
	"reflect"
	"strings"
)

func compareStringsIgnoreCase(target string) func(string) bool {
	target = strings.ToLower(target)
	return func(other string) bool {
		return strings.ToLower(other) == target
	}
}

// GetValueAtPath returns the value at the given path in source. The source must be a pointer.
// The value is returned as a reflect.Value in out.
func GetValueAtPath(source any, path ObjectPath, out *reflect.Value) error {
	value := reflect.ValueOf(source)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf(`cannot get value at path "%s": source is not a pointer`, path.String())
	}
	value = value.Elem()

	// Iterate over path elements
	for i, element := range path.elements {

		// Check if the value is zero or nil
		if !value.IsValid() {
			return fmt.Errorf(`cannot enter field "%s" of path %s at index %d: value is zero or nil`, element.name, path.String(), i)
		}

		// Dereference pointer
		switch value.Kind() {
		case reflect.Interface, reflect.Ptr:
			value = value.Elem()
		}

		// Check if we're working with a map or a struct
		switch value.Kind() {
		case reflect.Map:
			value = value.MapIndex(reflect.ValueOf(element.name)).Elem()
		case reflect.Struct:
			valueType := value.Type()
			field, ok := valueType.FieldByNameFunc(compareStringsIgnoreCase(element.name))
			if !ok {
				return fmt.Errorf(`cannot get value at path "%s": field "%s" not found in struct at path index %d`, path.String(), element.name, i)
			}
			value = value.FieldByIndex(field.Index)
		default:
			return fmt.Errorf(`cannot get value at path "%s": value at path index %d is neither a map nor struct`, path.String(), i)
		}
	}
	*out = value
	return nil
}

// AssignTypeAtPath assigns the given reflect.Type to the value at the given path in source.
// The source must be a pointer.
func AssignTypeAtPath(source any, path ObjectPath, newType reflect.Type) error {
	var value reflect.Value
	if err := GetValueAtPath(source, path, &value); err != nil {
		return err
	}

	// Set the new type
	value.Set(reflect.New(newType).Elem())
	return nil
}
