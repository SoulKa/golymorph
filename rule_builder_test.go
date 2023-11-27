package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
	"testing"
)

func rulesEqual(a Rule, b Rule) bool {
	return a.ValuePath.IsEqualTo(&b.ValuePath) &&
		a.NewType == b.NewType
}

func TestRuleBuilder(t *testing.T) {

	// Arrange
	valuePathString := "foo/bar"
	comparatorValue := "test"
	newType := reflect.TypeOf(int64(0))
	err, valuePath := objectpath.NewObjectPathFromString(valuePathString)
	if err != nil {
		t.Fatalf("error parsing input path [%s]: %s", valuePathString, err)
	}
	expectedRule := Rule{
		ValuePath:          *valuePath,
		ComparatorFunction: func(v any) bool { return v == comparatorValue },
		NewType:            newType,
	}

	// Act
	errors, rule := NewRuleBuilder().
		WhenValueAtPathString(valuePathString).
		IsEqualTo(comparatorValue).
		ThenAssignType(newType).
		Build()

	// Assert
	if len(errors) > 0 {
		for _, err := range errors {
			t.Log(err)
		}
		t.Fatalf("expected no errors, but got %d errors", len(errors))
	}
	if !rulesEqual(rule, expectedRule) {
		t.Fatalf("expected rule to be %+v, but got %+v", expectedRule, rule)
	}

}
