package golymorph_test

import (
	"encoding/json"
	"fmt"
	"github.com/SoulKa/golymorph"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// ExampleUnmarshalJSON demonstrates how to use the polymorpher to unmarshal a JSON into a struct with a polymorphic field.
func ExampleUnmarshalJSON() {

	// get a JSON that contains a payload with a type field that determines the type of the payload
	alertEventJson := `{ "timestamp": "2023-11-27T22:14:09+00:00", "payload": { "type": "alert", "message": "something is broken!" } }`

	// the parent type that contains the polymorphic payload
	type Event struct {
		Timestamp string
		Payload   any
	}

	// the polymorphic child types
	type AlertPayload struct {
		Type    string
		Message string
	}
	type PingPayload struct {
		Type string
		Ip   string
	}

	// define a mapping from the type value to the type of the payload
	typeMap := golymorph.TypeMap{
		"alert": reflect.TypeOf(AlertPayload{}),
		"ping":  reflect.TypeOf(PingPayload{}),
	}

	// create a TypeResolver that assigns the type of the payload based on the type field
	err, resolver := golymorph.NewPolymorphismBuilder().
		DefineTypeAt("payload").
		UsingTypeMap(typeMap).
		WithDiscriminatorAt("type").
		Build()
	if err != nil {
		panic(fmt.Sprintf("error building polymorpher: %s", err))
	}

	// create a new event
	var event Event
	if err := golymorph.UnmarshalJSON(resolver, []byte(alertEventJson), &event); err != nil {
		panic(fmt.Sprintf("error unmarshalling event: %s", err))
	}

	// continue to work with the event
	fmt.Printf("event: %+v\n", event)
	fmt.Printf("event payload: %T %+v\n", event.Payload, event.Payload.(AlertPayload))

	// Output:
	// event: {Timestamp:2023-11-27T22:14:09+00:00 Payload:{Type:alert Message:something is broken!}}
	// event payload: golymorph_test.AlertPayload {Type:alert Message:something is broken!}
}

// ExampleTypeMapPolymorphism_AssignTargetType demonstrates how to use the polymorpher to assign the
// type of a polymorphic field in an existing struct instance.
func ExampleTypeMapPolymorphism_AssignTargetType() {

	// get a JSON that contains a payload with a type field that determines the type of the payload
	alertEventJson := `{ "timestamp": "2023-11-27T22:14:09+00:00", "payload": { "type": "alert", "message": "something is broken!" } }`

	type Event struct {
		Timestamp string
		Payload   any
	}

	type AlertPayload struct {
		Type    string
		Message string
	}

	type PingPayload struct {
		Type string
		Ip   string
	}

	typeMap := golymorph.TypeMap{
		"alert": reflect.TypeOf(AlertPayload{}),
		"ping":  reflect.TypeOf(PingPayload{}),
	}

	// parse the JSON into a map
	var jsonMap map[string]any
	if err := json.Unmarshal([]byte(alertEventJson), &jsonMap); err != nil {
		panic(fmt.Sprintf("error unmarshalling JSON: %s", err))
	}

	// create a polymorpher that assigns the type of the payload based on the type field
	err, polymorpher := golymorph.NewPolymorphismBuilder().
		DefineTypeAt("payload").
		UsingTypeMap(typeMap).
		WithDiscriminatorAt("type").
		Build()
	if err != nil {
		panic(fmt.Sprintf("error building polymorpher: %s", err))
	}

	// create a new event
	var event Event
	if err := polymorpher.AssignTargetType(&jsonMap, &event); err != nil {
		panic(fmt.Sprintf("error assigning target type: %s", err))
	}

	// use mapstructure to unmarshal the payload into the event
	if err := mapstructure.Decode(jsonMap, &event); err != nil {
		panic(fmt.Sprintf("error decoding JSON map: %s", err))
	}

	// continue to work with the event
	fmt.Printf("event: %+v\n", event)
	fmt.Printf("event payload: %T %+v\n", event.Payload, event.Payload.(AlertPayload))

	// Output:
	// event: {Timestamp:2023-11-27T22:14:09+00:00 Payload:{Type:alert Message:something is broken!}}
	// event payload: golymorph_test.AlertPayload {Type:alert Message:something is broken!}
}
