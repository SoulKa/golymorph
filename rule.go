package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

// Rule is a rule for a polymorphism mapper.
type Rule struct {
	// ValuePath is the path to the value in the source to compare.
	ValuePath objectpath.ObjectPath
	// ComparatorFunction is the function to use to compare the value at ValuePath to.
	ComparatorFunction func(any) bool
	// NewType is the type to assign to the target if the rule matches.
	NewType reflect.Type
}

// Matches returns true if the source matches the rule.
func (r *Rule) Matches(source any) (error, bool) {
	var comparatorValue reflect.Value
	if err := objectpath.GetValueAtPath(source, r.ValuePath, &comparatorValue); err != nil {
		return err, false
	}
	return nil, r.ComparatorFunction(comparatorValue.Interface())
}
