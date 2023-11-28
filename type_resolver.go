package golymorph

import (
	"reflect"
)

type TypeMap map[any]reflect.Type

type TypeResolver interface {
	// AssignTargetType assigns the determined type to target based on the polymorphism rules. The matching rule with the
	// highest priority is used. If no rule matches, the target type is not changed. The source and target must be pointers.
	// If no matching type can be determined, an error.UnresolvedTypeError is returned.
	AssignTargetType(source any, target any) error
}
