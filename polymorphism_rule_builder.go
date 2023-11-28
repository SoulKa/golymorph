package golymorph

import (
	"errors"
	"github.com/SoulKa/golymorph/rules"
)

type polymorphismRuleBuilder struct {
	polymorphismBuilderBase
	rules []rules.Rule
}

func (b *polymorphismRuleBuilder) UsingRule(rule rules.Rule) polymorphismBuilderRuleAdder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *polymorphismRuleBuilder) Build() (error, TypeResolver) {
	if len(b.errors) > 0 {
		return errors.Join(b.errors...), nil
	}
	return nil, &RulePolymorphism{Polymorphism{b.targetPath}, b.rules}
}
