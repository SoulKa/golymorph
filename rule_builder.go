package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

// RuleBuilder is a builder for a polymorphism rule.
type RuleBuilder struct {
	errors         []error
	valuePath      objectpath.ObjectPath
	comparatorFunc func(any) bool
	newType        reflect.Type
}

type RuleBuilderBase interface {
	WhenValueAt(valuePath objectpath.ObjectPath) RuleBuilderConditionSetter
	WhenValueAtPathString(valuePath string) RuleBuilderConditionSetter
}

type RuleBuilderConditionSetter interface {
	IsEqualTo(value any) RuleBuilderTypeAssigner
	Matches(comparator func(any) bool) RuleBuilderTypeAssigner
}

type RuleBuilderTypeAssigner interface {
	ThenAssignType(newType reflect.Type) RuleBuilderFinalizer
}

type RuleBuilderFinalizer interface {
	Build() ([]error, Rule)
}

// NewRuleBuilder creates a new RuleBuilder. It enables a fluent interface for building a Rule.
func NewRuleBuilder() RuleBuilderBase {
	return &RuleBuilder{}
}

// WhenValueAt sets the path to the value in the source to compare.
func (b *RuleBuilder) WhenValueAt(valuePath objectpath.ObjectPath) RuleBuilderConditionSetter {
	b.valuePath = valuePath
	return b
}

// WhenValueAtPathString sets the path to the value in the source to compare.
func (b *RuleBuilder) WhenValueAtPathString(valuePath string) RuleBuilderConditionSetter {
	if err, path := objectpath.NewObjectPathFromString(valuePath); err != nil {
		b.appendError(err)
	} else {
		b.valuePath = *path
	}
	return b
}

// IsEqualTo sets the value to compare to.
func (b *RuleBuilder) IsEqualTo(value any) RuleBuilderTypeAssigner {
	b.comparatorFunc = func(v any) bool { return v == value }
	return b
}

// Matches sets the function to use to compare the value at ValuePath to.
func (b *RuleBuilder) Matches(comparator func(any) bool) RuleBuilderTypeAssigner {
	b.comparatorFunc = comparator
	return b
}

// ThenAssignType sets the type to assign to the target if the rule matches.
func (b *RuleBuilder) ThenAssignType(newType reflect.Type) RuleBuilderFinalizer {
	b.newType = newType
	return b
}

// Build builds the Rule and returns the errors encountered while building.
func (b *RuleBuilder) Build() ([]error, Rule) {
	return b.errors, Rule{
		b.valuePath,
		b.comparatorFunc,
		b.newType,
	}
}

func (b *RuleBuilder) appendError(err error) {
	b.errors = append(b.errors, err)
}
