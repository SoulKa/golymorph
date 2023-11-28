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
	// WhenValueAt sets the path to the value in the source to compare.
	WhenValueAt(valuePath string) ruleBuilderConditionSetter
}

type ruleBuilderConditionSetter interface {
	// IsEqualTo sets the value to compare to.
	IsEqualTo(value any) ruleBuilderTypeAssigner

	// Matches sets the function to use to compare the value at ValuePath to.
	Matches(comparator func(any) bool) ruleBuilderTypeAssigner
}

type ruleBuilderTypeAssigner interface {
	// ThenAssignType sets the type to assign to the target if the rule matches.
	ThenAssignType(newType reflect.Type) ruleBuilderFinalizer
}

type ruleBuilderFinalizer interface {
	// Build builds the Rule and returns the errors encountered while building.
	Build() ([]error, Rule)
}

// NewRuleBuilder creates a new ruleBuilder. It enables a fluent interface for building a Rule.
func NewRuleBuilder() ruleBuilderBase {
	return &ruleBuilder{}
}

func (b *ruleBuilder) WhenValueAt(valuePath string) ruleBuilderConditionSetter {
	if err, path := objectpath.NewObjectPathFromString(valuePath); err != nil {
		b.appendError(err)
	} else {
		b.valuePath = *path
	}
	return b
}

func (b *ruleBuilder) IsEqualTo(value any) ruleBuilderTypeAssigner {
	b.comparatorFunc = func(v any) bool { return v == value }
	return b
}

func (b *ruleBuilder) Matches(comparator func(any) bool) ruleBuilderTypeAssigner {
	b.comparatorFunc = comparator
	return b
}

func (b *ruleBuilder) ThenAssignType(newType reflect.Type) ruleBuilderFinalizer {
	b.newType = newType
	return b
}

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
