package golymorph

import (
	"errors"
	golimorphError "github.com/SoulKa/golymorph/error"
	"github.com/SoulKa/golymorph/objectpath"
)

// RulePolymorphism is a mapper that assigns a target type based on the given Rules
type RulePolymorphism struct {
	Polymorphism

	// Rules is a list of Rules to apply. The first rule that matches is used to determine the target type.
	Rules []Rule
}

func (p *RulePolymorphism) AssignTargetType(source any, target any) error {

	// check for each rule if it matches and assign type if it does
	for _, rule := range p.Rules {
		if err, matches := rule.Matches(source); err != nil {
			return errors.Join(errors.New("error applying rule"), err)
		} else if matches {
			if err := objectpath.AssignTypeAtPath(target, p.TargetPath, rule.NewType); err != nil {
				return errors.Join(errors.New("error assigning type to target"), err)
			}
			return nil
		}
	}

	// no rule matched
	return &golimorphError.UnresolvedTypeError{
		Err:        errors.New("no rule matched"),
		TargetPath: p.TargetPath.String(),
	}
}
