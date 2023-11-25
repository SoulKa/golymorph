package golymorph

import (
	"fmt"
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

type ComparatorType int

const (
	ComparatorTypeEquality ComparatorType = iota
)

// Rule is a rule for a polymorphism mapper.
type Rule struct {
	// ValuePath is the path to the value in the source to compare.
	ValuePath objectpath.ObjectPath
	// ComparatorType is the type of comparison to perform.
	ComparatorType ComparatorType
	// ComparatorValue is the value to compare against.
	ComparatorValue any
	// NewType is the type to assign to the target if the rule matches.
	NewType reflect.Type
}

// Matches returns true if the source matches the rule.
func (r *Rule) Matches(source any) (error, bool) {
	var sourceValue reflect.Value
	if err := objectpath.GetValueAtPath(source, r.ValuePath, &sourceValue); err != nil {
		return err, false
	}

	switch r.ComparatorType {
	case ComparatorTypeEquality:
		return nil, sourceValue.Interface() == r.ComparatorValue
	default:
		return fmt.Errorf("unknown comparator type %d", r.ComparatorType), false
	}
}
