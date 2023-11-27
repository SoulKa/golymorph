package golymorph

import (
	"encoding/json"
	"github.com/SoulKa/golymorph/objectpath"
	"github.com/mitchellh/mapstructure"
)

type Polymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath
}

// UnmarshalJSON unmarshals the given JSON data into the given output object using the given TypeResolver.
func UnmarshalJSON(resolver TypeResolver, data []byte, output any) (error, bool) {

	// parse JSON
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err, false
	}

	// resolve polymorphism
	if err, b := Decode(resolver, jsonMap, output); err != nil {
		return err, false
	} else if !b {
		return nil, false
	}

	// success
	return nil, true
}

// Decode decodes the given source map into the given output object using the given TypeResolver and mapstructure.
func Decode(resolver TypeResolver, source map[string]any, output any) (error, bool) {

	// create a new event
	err, assigned := resolver.AssignTargetType(&source, output)
	if err != nil {
		return err, false
	} else if !assigned {
		return nil, false
	}

	// use mapstructure to unmarshal the payload into the event
	if err := mapstructure.Decode(source, output); err != nil {
		return err, false
	}

	// success
	return nil, true
}
