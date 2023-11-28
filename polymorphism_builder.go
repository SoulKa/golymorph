package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
	"github.com/SoulKa/golymorph/rules"
	"strings"
)

type polymorphismBuilderBase struct {
	targetPath objectpath.ObjectPath
	errors     []error
}

type polymorphismBuilderEmpty interface {
	// DefineTypeAt defines the target path of the polymorphism. This is the path where the polymorphism
	// will be applied, i.e. where the new type is set. For valid paths see objectpath.NewObjectPathFromString.
	DefineTypeAt(targetPath string) polymorphismBuilderStrategySelector
}

type polymorphismBuilderStrategySelector interface {
	// UsingRule defines a rule that is used to determine the new type. The rules are applied in the
	// order they are defined. The first rule that matches is used to determine the new type.
	UsingRule(rule rules.Rule) polymorphismBuilderRuleAdder

	// UsingTypeMap defines a type map that is used to determine the new type. The type map is applied
	UsingTypeMap(typeMap TypeMap) polymorphismBuilderDiscriminatorKeyDefiner
}

type polymorphismBuilderRuleAdder interface {
	// UsingRule defines a rule that is used to determine the new type. The rules are applied in the
	// order they are defined. The first rule that matches is used to determine the new type.
	UsingRule(rule rules.Rule) polymorphismBuilderRuleAdder

	// Build creates a new TypeResolver that can be used to resolve a polymorphic type.
	Build() (error, TypeResolver)
}

type polymorphismBuilderDiscriminatorKeyDefiner interface {
	// WithDiscriminatorAt defines the path to the discriminator key. The discriminator key is used to
	// determine the new type. The value of the discriminator key is used to lookup the new type in the
	// type map.
	WithDiscriminatorAt(discriminatorKey string) polymorphismBuilderFinalizer
}

type polymorphismBuilderFinalizer interface {
	// Build creates a new TypeResolver that can be used to resolve a polymorphic type.
	Build() (error, TypeResolver)
}

// NewPolymorphismBuilder creates a new polymorphism builder that is used in a human readable way to create a polymorphism.
// It only allows a valid combination of rules and type maps.
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

func (b *polymorphismBuilderBase) UsingRule(rule rules.Rule) polymorphismBuilderRuleAdder {
	return &polymorphismRuleBuilder{
		polymorphismBuilderBase: *b,
		rules:                   []rules.Rule{rule},
	}
}

func (b *polymorphismBuilderBase) UsingTypeMap(typeMap TypeMap) polymorphismBuilderDiscriminatorKeyDefiner {
	return &polymorphismTypeMapBuilder{
		polymorphismBuilderBase: *b,
		typeMap:                 typeMap,
	}
}
