# golymorph [![Pipeline Tests Status](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml/badge.svg)](https://github.com/SoulKa/golymorph/actions/workflows/go-test.yaml) [![Godoc](https://godoc.org/github.com/SoulKa/golymorph?status.svg)](https://godoc.org/github.com/SoulKa/golymorph)
The golymorph module enables resolving polymorphic typing at runtime. It's usually used in conjunction with
JSON parsing and the `mapstructure` module. In fact, this module takes the use case of `mapstructure` and takes it
a step further by allowing the user to define a custom type resolver function.

## Installation

Standard `go get`:

```shell
go get github.com/SoulKa/golymorph
```

## Usage & Example

For usage and examples see the [Godoc](http://godoc.org/github.com/SoulKa/golymorph).

```go
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

typeMap := TypeMap{
    "alert": reflect.TypeOf(AlertPayload{}),
    "ping":  reflect.TypeOf(PingPayload{}),
}

// parse the JSON into a map
var jsonMap map[string]any
if err := json.Unmarshal([]byte(alertEventJson), &jsonMap); err != nil {
    t.Fatalf("error unmarshalling JSON: %s", err)
}

// create a polymorpher that assigns the type of the payload based on the type field
err, polymorpher := NewPolymorphismBuilder().
    DefineTypeAt("payload").
    UsingTypeMap(typeMap).
    WithDiscriminatorAt("type").
    Build()
if err != nil {
    t.Fatalf("error building polymorpher: %s", err)
}

// create a new event
var event Event
err, assigned := polymorpher.AssignTargetType(&jsonMap, &event)
if err != nil {
    t.Fatalf("error assigning target type: %s", err)
} else if !assigned {
    t.Fatalf("no type assigned")
}

// use mapstructure to unmarshal the payload into the event
if err := mapstructure.Decode(jsonMap, &event); err != nil {
    t.Fatalf("error decoding JSON map: %s", err)
}

// continue to work with the event
fmt.Printf("event: %+v\n", event)
fmt.Printf("event payload: %T %+v\n", event.Payload, event.Payload.(AlertPayload))
```