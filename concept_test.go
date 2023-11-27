package golymorph

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"testing"
)

func TestReadingMapByField(t *testing.T) {
	input := map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"test": 123,
			},
		},
	}

	t.Logf("input: %v\n", input)
	value := reflect.ValueOf(input)
	foo := value.MapIndex(reflect.ValueOf("foo")).Elem()
	t.Logf("foo: %v\n", foo)
	bar := foo.MapIndex(reflect.ValueOf("bar")).Elem()
	t.Logf("bar: %v\n", bar)
	test := bar.MapIndex(reflect.ValueOf("test")).Elem()
	t.Logf("test: %v\n", test)
}

func TestMapstructure(t *testing.T) {
	type Bar struct {
		Test int
	}

	type Foo struct {
		Bar  any
		Test int
	}

	fooVal := reflect.ValueOf(Foo{})
	fooType := fooVal.Type()
	for i := 0; i < fooType.NumField(); i++ {
		field := fooType.Field(i)
		t.Logf("fooVal.Field(%d): %s=%v\n", i, field.Name, fooVal.Field(i))
	}

	input := map[string]any{
		"Bar": map[string]int{
			"Test": 123,
		},
		"Test": 456,
	}

	t.Logf("input: %v\n", input)

	// decode
	var output Foo
	if err := mapstructure.Decode(input, &output); err != nil {
		t.Errorf("error decoding: %s", err)
	}
	t.Logf("output: %T%+v\n", output, output)
}

func TestWritingStructByFields(t *testing.T) {
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

	animal := Animal{
		Name:      "horsey",
		Specifics: map[string]any{"type": "horse", "shoes": 4},
	}

	value := reflect.ValueOf(&animal).Elem()
	name := value.FieldByName("Name")
	t.Logf("name: %v\n", name)
	value = value.FieldByName("Specifics")
	t.Logf("specifics: %v\n", value)
	value.Set(reflect.ValueOf(Horse{4}))
	t.Logf("animal: %+v\n", animal)

	if reflect.TypeOf(animal.Specifics) != reflect.TypeOf(Horse{}) {
		t.Fatalf("expected animal.Specifics to be of type Horse, but got %T", animal.Specifics)
	}
}

func TestWritingStructByFieldsInMethod(t *testing.T) {
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

	animal := Animal{
		Name:      "horsey",
		Specifics: map[string]any{"type": "horse", "shoes": 4},
	}

	SetField(&animal)
	t.Logf("animal: %+v\n", animal)

	if reflect.TypeOf(animal.Specifics).Name() != reflect.TypeOf(Horse{4}).Name() {
		t.Fatalf("expected animal.Specifics to be of type %+v, but got %+v", reflect.TypeOf(Horse{4}), reflect.TypeOf(animal.Specifics))
	}
}

func SetField(animalPtr any) {
	value := reflect.ValueOf(animalPtr).Elem()
	name := value.FieldByName("Name")
	fmt.Printf("name: %v\n", name)
	value = value.FieldByName("Specifics")
	fmt.Printf("specifics: %v\n", value)
	value.Set(reflect.ValueOf(Horse{4}))
	fmt.Printf("animal: %+v\n", reflect.ValueOf(animalPtr).Elem().Interface())
}

func TestCompareAny(t *testing.T) {
	a := any(1)
	b := any(1)
	if a != b {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}

func TestAnyMap(t *testing.T) {
	typeVal := reflect.TypeOf(int64(0))
	anyMap := TypeMap{
		"foo": typeVal,
	}
	if anyMap["foo"] != typeVal {
		t.Fatalf(`expected anyMap["foo"] to equal %+v`, typeVal)
	}
}
