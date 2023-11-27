package golymorph

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"testing"
)

type Animal struct {
	Name      string
	Specifics any
}

type Horse struct {
	Shoes int
}

type Duck struct {
	Feathers int
}

var animalTypeMap = TypeMap{
	"horse": reflect.TypeOf(Horse{}),
	"duck":  reflect.TypeOf(Duck{}),
}

type TestCase struct {
	inputJson string
	output    Animal
}

var testCases = []TestCase{
	{`{ "name": "horsey", "specifics": { "type": "horse", "shoes": 4 } }`, Animal{"horsey", Horse{4}}},
	{`{ "name": "ducky", "specifics": { "type": "duck", "feathers": 1000 } }`, Animal{"ducky", Duck{1000}}},
}

func TestPolymorphism_AssignTargetType(t *testing.T) {

	// Arrange
	errs, polymorphism := NewPolymorphismBuilder().
		DefineTypeAt("specifics").
		UsingTypeMap(animalTypeMap).
		WithDiscriminatorAt("type").
		Build()
	if HasErrors(t, errs) {
		t.Fatalf("error creating polymorphism: %s", errs)
	}
	t.Logf("polymorphism: %+v\n", polymorphism)

	for _, tc := range testCases {

		// parse JSON
		var actualAnimalJson any
		if err := json.Unmarshal([]byte(tc.inputJson), &actualAnimalJson); err != nil {
			t.Fatalf("error unmarshalling horse: %s", err)
		}
		t.Logf("actualAnimalJson: %+v\n", actualAnimalJson)

		// Act
		var actualAnimal Animal
		if err, applied := polymorphism.AssignTargetType(&actualAnimalJson, &actualAnimal); err != nil {
			t.Fatalf("error assigning target type to horse: %s", err)
		} else if !applied {
			t.Fatalf("expected polymorphism to be applied")
		}
		t.Logf("actualAnimal: %+v\n", actualAnimal)

		// map map to struct
		if err := mapstructure.Decode(actualAnimalJson, &actualAnimal); err != nil {
			t.Fatalf("error decoding animal: %s\n", err)
		}
		t.Logf("actualAnimal: %+v\n", actualAnimal)

		// Assert
		if !reflect.DeepEqual(actualAnimal, tc.output) {
			t.Fatalf("expected horse to be %+v, but got %+v", tc.output, actualAnimal)
		}
	}
}