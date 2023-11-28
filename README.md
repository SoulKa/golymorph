# golymorph [![Pipeline Tests Status](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml/badge.svg)](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml) [![Godoc](https://godoc.org/github.com/SoulKa/golymorph?status.svg)](https://godoc.org/github.com/SoulKa/golymorph)

The golymorph module enables resolving polymorphic typing at runtime. It's usually used in
conjunction with JSON parsing and the `mapstructure` module. In fact, this module takes the use case
of `mapstructure` a step further by allowing the user to define a custom type resolver.

## Installation

Standard `go get`:

```shell
go get github.com/SoulKa/golymorph
```

## Docs

The docs are hosted on [Godoc](http://godoc.org/github.com/SoulKa/golymorph).

## Use Case

Use this module to resolve polymorphic types at runtime. An example would be a struct that contains
a payload field which can be of different struct types not known at compile time:

```go
// the parent type that contains the polymorphic payload
type Event struct {
	Timestamp string
	Payload   any // <-- AlertPayload or PingPayload?
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
```

If, for example, you have a JSON that you decode into the `Event` struct, it is cumbersome to parse the JSON into a map, look into the `type` field of the `payload` and after that select and parse the map into the correct type at the `payload` field of the `event` struct.
golymorph does exactly this:

1. Look at the value of a defined field anywhere in the given `map` or JSON
2. Find the correct type using the given pairs of `value ==> reflect.Type`
3. Assign the correct type at the given position anywhere in the "parent" struct
4. Fully decode the given JSON or `map`, now with a concrete struct type as `payload`

## Example

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
