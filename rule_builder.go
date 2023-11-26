package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

type RuleBuilder struct {
	errors          []error
	valuePath       objectpath.ObjectPath
	comparatorType  ComparatorType
	comparatorValue any
	newType         reflect.Type
}

type RuleBuilderBase interface {
	WhenValueAt(valuePath objectpath.ObjectPath) RuleBuilderPathSet
	WhenValueAtPathString(valuePath string) RuleBuilderPathSet
}

type RuleBuilderPathSet interface {
	IsEqualTo(value any) RuleBuilderConditionSet
	Matches(comparator func(any) bool) RuleBuilderConditionSet
}

type RuleBuilderConditionSet interface {
	ThenAssignType(newType reflect.Type) RuleBuilderNewTypeSet
}

type RuleBuilderNewTypeSet interface {
	BuildRule() Rule
}

func NewRuleBuilder() RuleBuilderBase {
	return &RuleBuilder{}
}

func (b *RuleBuilder) WhenValueAt(valuePath objectpath.ObjectPath) RuleBuilderPathSet {
	b.valuePath = valuePath
	return b
}

func (b *RuleBuilder) WhenValueAtPathString(valuePath string) RuleBuilderPathSet {
	if err, path := objectpath.NewObjectPathFromString(valuePath); err != nil {
		b.appendError(err)
	} else {
		b.valuePath = *path
	}
	return b
}

func (b *RuleBuilder) IsEqualTo(value any) RuleBuilderConditionSet {
	b.comparatorType = ComparatorTypeEquality
	b.comparatorValue = value
	return b
}

func (b *RuleBuilder) Matches(comparator func(any) bool) RuleBuilderConditionSet {
	b.comparatorType = ComparatorTypeFunction
	b.comparatorValue = comparator
	return b
}

func (b *RuleBuilder) ThenAssignType(newType reflect.Type) RuleBuilderNewTypeSet {
	b.newType = newType
	return b
}

func (b *RuleBuilder) BuildRule() Rule {
	return Rule{
		b.valuePath,
		b.comparatorType,
		b.comparatorValue,
		b.newType,
	}
}

func (b *RuleBuilder) Errors() []error {
	return b.errors
}

func (b *RuleBuilder) appendError(err error) {
	b.errors = append(b.errors, err)
}
