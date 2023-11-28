package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"strings"
)

type polymorphismBuilderBase struct {
	targetPath objectpath.ObjectPath
	errors     []error
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

func (b *polymorphismBuilderBase) UsingTypeMap(typeMap TypeMap) polymorphismBuilderDiscriminatorKeyDefiner {
	return &polymorphismTypeMapBuilder{
		polymorphismBuilderBase: *b,
		typeMap:                 typeMap,
	}
}
