package golymorph

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// TypeMap is a map of values to types.
type TypeMap map[any]reflect.Type

// TypeResolver is an interface that can resolve the type of a target based on the values of a source.
type TypeResolver interface {
	// AssignTargetType assigns the determined type to target based on the polymorphism rules. The matching rule with the
	// highest priority is used. If no rule matches, the target type is not changed. The source and target must be pointers.
	// If no matching type can be determined, an error.UnresolvedTypeError is returned.
	AssignTargetType(source any, target any) error
}

// UnmarshalJSON unmarshals the given JSON data into the given output object using the given TypeResolver.
func UnmarshalJSON(resolver TypeResolver, data []byte, output any) error {

	// parse JSON
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	// resolve polymorphism
	if err := Decode(resolver, jsonMap, output); err != nil {
		return err
	}

	// success
	return nil
}

// Decode the given source map into the given output object using the given TypeResolver and mapstructure.
// The output object must be a pointer.
func Decode(resolver TypeResolver, source map[string]any, output any) error {

	// create a new event
	if err := resolver.AssignTargetType(&source, output); err != nil {
		return err
	}

	// use mapstructure to unmarshal the payload into the event
	if err := mapstructure.Decode(source, output); err != nil {
		return err
	}

	// success
	return nil
}
