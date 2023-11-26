package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
	"testing"
)

func rulesEqual(a Rule, b Rule) bool {
	return a.ValuePath.IsEqualTo(&b.ValuePath) &&
		a.ComparatorType == b.ComparatorType &&
		a.ComparatorValue == b.ComparatorValue &&
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
		ValuePath:       *valuePath,
		ComparatorType:  ComparatorTypeEquality,
		ComparatorValue: comparatorValue,
		NewType:         newType,
	}

	// Act
	rule := NewRuleBuilder().
		WhenValueAtPathString(valuePathString).
		IsEqualTo(comparatorValue).
		ThenAssignType(newType).
		BuildRule()

	// Assert
	if !rulesEqual(rule, expectedRule) {
		t.Fatalf("expected rule to be %+v, but got %+v", expectedRule, rule)
	}

}
