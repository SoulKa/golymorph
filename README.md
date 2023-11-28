# golymorph [![Pipeline Tests Status](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml/badge.svg)](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml) [![Godoc](https://godoc.org/github.com/SoulKa/golymorph?status.svg)](https://godoc.org/github.com/SoulKa/golymorph)

The golymorph module enables resolving polymorphic typing at runtime. It's usually used in
conjunction with
JSON parsing and the `mapstructure` module. In fact, this module takes the use case
of `mapstructure` and takes it
a step further by allowing the user to define a custom type resolver function.

## Installation

Standard `go get`:

```shell
go get github.com/SoulKa/golymorph
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/SoulKa/golymorph).

```go
package main

import (
	"fmt"
	"github.com/SoulKa/golymorph"
	"reflect"
)

func main() {
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
}
```