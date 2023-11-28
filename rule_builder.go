package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

type ruleBuilder struct {
	errors         []error
	valuePath      objectpath.ObjectPath
	comparatorFunc func(any) bool
	newType        reflect.Type
}

type ruleBuilderBase interface {
	WhenValueAt(valuePath objectpath.ObjectPath) ruleBuilderConditionSetter
	WhenValueAtPathString(valuePath string) ruleBuilderConditionSetter
}

type ruleBuilderConditionSetter interface {
	IsEqualTo(value any) ruleBuilderTypeAssigner
	Matches(comparator func(any) bool) ruleBuilderTypeAssigner
}

type ruleBuilderTypeAssigner interface {
	ThenAssignType(newType reflect.Type) ruleBuilderFinalizer
}

type ruleBuilderFinalizer interface {
	Build() ([]error, Rule)
}

// NewRuleBuilder creates a new ruleBuilder. It enables a fluent interface for building a Rule.
func NewRuleBuilder() ruleBuilderBase {
	return &ruleBuilder{}
}

// WhenValueAt sets the path to the value in the source to compare.
func (b *ruleBuilder) WhenValueAt(valuePath objectpath.ObjectPath) ruleBuilderConditionSetter {
	b.valuePath = valuePath
	return b
}

// WhenValueAtPathString sets the path to the value in the source to compare.
func (b *ruleBuilder) WhenValueAtPathString(valuePath string) ruleBuilderConditionSetter {
	if err, path := objectpath.NewObjectPathFromString(valuePath); err != nil {
		b.appendError(err)
	} else {
		b.valuePath = *path
	}
	return b
}

// IsEqualTo sets the value to compare to.
func (b *ruleBuilder) IsEqualTo(value any) ruleBuilderTypeAssigner {
	b.comparatorFunc = func(v any) bool { return v == value }
	return b
}

// Matches sets the function to use to compare the value at ValuePath to.
func (b *ruleBuilder) Matches(comparator func(any) bool) ruleBuilderTypeAssigner {
	b.comparatorFunc = comparator
	return b
}

// ThenAssignType sets the type to assign to the target if the rule matches.
func (b *ruleBuilder) ThenAssignType(newType reflect.Type) ruleBuilderFinalizer {
	b.newType = newType
	return b
}

// Build builds the Rule and returns the errors encountered while building.
func (b *ruleBuilder) Build() ([]error, Rule) {
	return b.errors, Rule{
		b.valuePath,
		b.comparatorFunc,
		b.newType,
	}
}

func (b *ruleBuilder) appendError(err error) {
	b.errors = append(b.errors, err)
}
