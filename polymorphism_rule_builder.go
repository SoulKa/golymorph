package golymorph

import (
	"errors"
)

type polymorphismRuleBuilder struct {
	polymorphismBuilderBase
	rules []Rule
}

func (b *polymorphismRuleBuilder) UsingRule(rule Rule) polymorphismBuilderRuleAdder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *polymorphismRuleBuilder) Build() (error, TypeResolver) {
	if len(b.errors) > 0 {
		return errors.Join(b.errors...), nil
	}
	return nil, &RulePolymorphism{
		Polymorphism{
			TargetPath: b.targetPath},
		b.rules}
}
