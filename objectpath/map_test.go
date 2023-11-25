package objectpath

import (
	"reflect"
	"testing"
)

type TestCase struct {
	inputObject any
	inputPath   string
	output      any
}

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

func TestGetValueAtPath(t *testing.T) {
	var testCases = []TestCase{
		{map[string]any{"foo": map[string]any{"bar": map[string]any{"test": int64(123)}}}, "foo/bar/test", int64(123)},
	}

	for _, tc := range testCases {

		// Arrange
		var outVal reflect.Value
		err, inputPath := NewObjectPathFromString(tc.inputPath)
		if err != nil {
			t.Fatalf("error parsing input path %s: %s", tc.inputPath, err)
		}
		input := tc.inputObject

		// Act
		if err := GetValueAtPath(&input, *inputPath, &outVal); err != nil {
			t.Fatalf("error getting value at path %s: %s", tc.inputPath, err)
		}

		// Assert
		if outVal.Int() != tc.output {
			t.Fatalf("expected output to be %v, but got %v", tc.output, outVal)
		}
	}
}

func TestAssignTypeAtPath(t *testing.T) {
	var testCases = []TestCase{
		{Animal{Name: "horse", Specifics: map[string]any{}}, "Specifics", Horse{}},
		{Animal{Name: "duck", Specifics: map[string]any{}}, "Specifics", Duck{}},
	}

	for _, tc := range testCases {

		// Arrange
		err, inputPath := NewObjectPathFromString(tc.inputPath)
		if err != nil {
			t.Fatalf("error parsing input path %s: %s", tc.inputPath, err)
		}
		newType := reflect.TypeOf(tc.output)
		animal := tc.inputObject.(Animal)

		// Act
		if err := AssignTypeAtPath(&animal, *inputPath, newType); err != nil {
			t.Fatalf("error assigning type at path %s: %s", tc.inputPath, err)
		}
		t.Logf("animal: %+v\n", animal)

		// Assert
		outputType := reflect.TypeOf(animal.Specifics)
		if outputType != newType {
			t.Fatalf("expected output to be %v, but got %v", newType, outputType)
		}
	}
}
