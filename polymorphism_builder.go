package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"strings"
)

type PolymorphismBuilderBase struct {
	targetPath objectpath.ObjectPath
	errors     []error
}

type PolymorphismTypeMapBuilder struct {
	PolymorphismBuilderBase
	typeMap           TypeMap
	discriminatorPath objectpath.ObjectPath
}

type PolymorphismRuleBuilder struct {
	PolymorphismBuilderBase
	rules []Rule
}

type PolymorphismBuilderEmpty interface {
	DefineTypeAt(targetPath string) PolymorphismBuilderStrategySelector
}

type PolymorphismBuilderStrategySelector interface {
	UsingRule(rule Rule) PolymorphismBuilderRuleAdder
	UsingTypeMap(typeMap TypeMap) PolymorphismBuilderDiscriminatorKeyDefiner
}

type PolymorphismBuilderRuleAdder interface {
	UsingRule(rule Rule) PolymorphismBuilderRuleAdder
	Build() ([]error, Polymorpher)
}

type PolymorphismBuilderDiscriminatorKeyDefiner interface {
	WithDiscriminatorAt(discriminatorKey string) PolymorphismBuilderFinalizer
}

type PolymorphismBuilderFinalizer interface {
	Build() ([]error, Polymorpher)
}

func NewPolymorphismBuilder() PolymorphismBuilderEmpty {
	return &PolymorphismBuilderBase{*objectpath.NewSelfReferencePath(), []error{}}
}

func (b *PolymorphismBuilderBase) DefineTypeAt(targetPath string) PolymorphismBuilderStrategySelector {
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

func (b *PolymorphismBuilderBase) UsingRule(rule Rule) PolymorphismBuilderRuleAdder {
	return &PolymorphismRuleBuilder{
		PolymorphismBuilderBase: *b,
		rules:                   []Rule{rule},
	}
}

func (b *PolymorphismRuleBuilder) UsingRule(rule Rule) PolymorphismBuilderRuleAdder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *PolymorphismBuilderBase) UsingTypeMap(typeMap TypeMap) PolymorphismBuilderDiscriminatorKeyDefiner {
	return &PolymorphismTypeMapBuilder{
		PolymorphismBuilderBase: *b,
		typeMap:                 typeMap,
	}
}

func (b *PolymorphismTypeMapBuilder) WithDiscriminatorAt(discriminatorKey string) PolymorphismBuilderFinalizer {
	if err, path := objectpath.NewObjectPathFromString(discriminatorKey); err != nil {
		b.errors = append(b.errors, err)
	} else if err := path.ToAbsolutePath(&b.targetPath); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.discriminatorPath = *path
	}
	return b
}

func (b *PolymorphismRuleBuilder) Build() ([]error, Polymorpher) {
	if len(b.errors) > 0 {
		return b.errors, nil
	}
	return b.errors, &RulePolymorphism{b.targetPath, b.rules}
}

func (b *PolymorphismTypeMapBuilder) Build() ([]error, Polymorpher) {
	if len(b.errors) > 0 {
		return b.errors, nil
	}
	return b.errors, &TypeMapPolymorphism{b.targetPath, b.discriminatorPath, b.typeMap}
}
