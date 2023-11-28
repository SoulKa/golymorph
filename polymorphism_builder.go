package golymorph

import (
	"errors"
	"github.com/SoulKa/golymorph/objectpath"
	"strings"
)

type polymorphismBuilderBase struct {
	targetPath objectpath.ObjectPath
	errors     []error
}

type polymorphismTypeMapBuilder struct {
	polymorphismBuilderBase
	typeMap           TypeMap
	discriminatorPath objectpath.ObjectPath
}

type polymorphismRuleBuilder struct {
	polymorphismBuilderBase
	rules []Rule
}

type polymorphismBuilderEmpty interface {
	DefineTypeAt(targetPath string) polymorphismBuilderStrategySelector
}

type polymorphismBuilderStrategySelector interface {
	UsingRule(rule Rule) polymorphismBuilderRuleAdder
	UsingTypeMap(typeMap TypeMap) polymorphismBuilderDiscriminatorKeyDefiner
}

type polymorphismBuilderRuleAdder interface {
	UsingRule(rule Rule) polymorphismBuilderRuleAdder
	Build() (error, TypeResolver)
}

type polymorphismBuilderDiscriminatorKeyDefiner interface {
	WithDiscriminatorAt(discriminatorKey string) polymorphismBuilderFinalizer
}

type polymorphismBuilderFinalizer interface {
	Build() (error, TypeResolver)
}

func NewPolymorphismBuilder() polymorphismBuilderEmpty {
	return &polymorphismBuilderBase{*objectpath.NewSelfReferencePath(), []error{}}
}

func (b *polymorphismBuilderBase) DefineTypeAt(targetPath string) polymorphismBuilderStrategySelector {
	// make target path absolute
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath
	}

	// parse target path
	if err, path := objectpath.NewObjectPathFromString(targetPath); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.targetPath = *path
	}
	return b
}

func (b *polymorphismBuilderBase) UsingRule(rule Rule) polymorphismBuilderRuleAdder {
	return &polymorphismRuleBuilder{
		polymorphismBuilderBase: *b,
		rules:                   []Rule{rule},
	}
}

func (b *polymorphismRuleBuilder) UsingRule(rule Rule) polymorphismBuilderRuleAdder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *polymorphismBuilderBase) UsingTypeMap(typeMap TypeMap) polymorphismBuilderDiscriminatorKeyDefiner {
	return &polymorphismTypeMapBuilder{
		polymorphismBuilderBase: *b,
		typeMap:                 typeMap,
	}
}

func (b *polymorphismTypeMapBuilder) WithDiscriminatorAt(discriminatorKey string) polymorphismBuilderFinalizer {
	if err, path := objectpath.NewObjectPathFromString(discriminatorKey); err != nil {
		b.errors = append(b.errors, err)
	} else if err := path.ToAbsolutePath(&b.targetPath); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.discriminatorPath = *path
	}
	return b
}

func (b *polymorphismRuleBuilder) Build() (error, TypeResolver) {
	if len(b.errors) > 0 {
		return errors.Join(b.errors...), nil
	}
	return nil, &RulePolymorphism{Polymorphism{b.targetPath}, b.rules}
}

func (b *polymorphismTypeMapBuilder) Build() (error, TypeResolver) {
	if len(b.errors) > 0 {
		return errors.Join(b.errors...), nil
	}
	return nil, &TypeMapPolymorphism{Polymorphism{b.targetPath}, b.discriminatorPath, b.typeMap}
}
