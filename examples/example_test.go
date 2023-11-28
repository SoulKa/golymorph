package examples

import (
	"encoding/json"
	"fmt"
	"github.com/SoulKa/golymorph"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"testing"
)

func TestBasicPolymorphismFromJson(t *testing.T) {

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
		t.Fatalf("error building polymorpher: %s", err)
	}

	// create a new event
	var event Event
	if err := golymorph.UnmarshalJSON(resolver, []byte(alertEventJson), &event); err != nil {
		t.Fatalf("error unmarshalling event: %s", err)
	}

	// continue to work with the event
	fmt.Printf("event: %+v\n", event)
	fmt.Printf("event payload: %T %+v\n", event.Payload, event.Payload.(AlertPayload))
}

func TestBasicPolymorphismWithManualParsing(t *testing.T) {

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
		t.Fatalf("error unmarshalling JSON: %s", err)
	}

	// create a polymorpher that assigns the type of the payload based on the type field
	err, polymorpher := golymorph.NewPolymorphismBuilder().
		DefineTypeAt("payload").
		UsingTypeMap(typeMap).
		WithDiscriminatorAt("type").
		Build()
	if err != nil {
		t.Fatalf("error building polymorpher: %s", err)
	}

	// create a new event
	var event Event
	if err := polymorpher.AssignTargetType(&jsonMap, &event); err != nil {
		t.Fatalf("error assigning target type: %s", err)
	}

	// use mapstructure to unmarshal the payload into the event
	if err := mapstructure.Decode(jsonMap, &event); err != nil {
		t.Fatalf("error decoding JSON map: %s", err)
	}

	// continue to work with the event
	fmt.Printf("event: %+v\n", event)
	fmt.Printf("event payload: %T %+v\n", event.Payload, event.Payload.(AlertPayload))
}
