package golymorph

import (
	"reflect"
	"testing"
)

func TestPolymorphismBuilder_UsingRule(t *testing.T) {

	// Arrange
	errors, rule1 := NewRuleBuilder().
		WhenValueAt("foo/bar").
		IsEqualTo("test").
		ThenAssignType(reflect.TypeOf(int64(0))).
		Build()
	if HasErrors(t, errors) {
		t.Fatalf("expected no errors, but got %d errors", len(errors))
	}

	// Arrange
	errors, rule2 := NewRuleBuilder().
		WhenValueAt("foo/bar").
		IsEqualTo("test").
		ThenAssignType(reflect.TypeOf(int64(0))).
		Build()
	if HasErrors(t, errors) {
		t.Fatalf("expected no errors, but got %d errors", len(errors))
	}

	// Act
	err, polymorphism := NewPolymorphismBuilder().
		DefineTypeAt("foo/bar").
		UsingRule(rule1).
		UsingRule(rule2).
		Build()

	// Assert
	if err != nil {
		t.Fatalf("expected no errors, but got %s", err)
	} else if polymorphism == nil {
		t.Fatalf("expected polymorphism to not be nil")
	}
}

func TestPolymorphismBuilder_UsingTypeMap(t *testing.T) {

	// Arrange
	typeMap := TypeMap{
		"test": reflect.TypeOf(int64(0)),
	}

	// Act
	err, polymorphism := NewPolymorphismBuilder().
		DefineTypeAt("foo/bar").
		UsingTypeMap(typeMap).
		WithDiscriminatorAt("foo/bar/discriminator").
		Build()

	// Assert
	if err != nil {
		t.Fatalf("expected no errors, but got %s", err)
	} else if polymorphism == nil {
		t.Fatalf("expected polymorphism to not be nil")
	}
}

func HasErrors(t *testing.T, errors []error) bool {
	for _, err := range errors {
		t.Error(err)
	}
	return len(errors) > 0
}
