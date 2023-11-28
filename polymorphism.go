package golymorph

import (
	"encoding/json"
	"github.com/SoulKa/golymorph/objectpath"
	"github.com/mitchellh/mapstructure"
)

type Polymorphism struct {
	// TargetPath is the path to the object to assign the new type to
	TargetPath objectpath.ObjectPath
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

// Decode decodes the given source map into the given output object using the given TypeResolver and mapstructure.
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
