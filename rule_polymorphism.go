package golymorph

import (
	"errors"
	"github.com/SoulKa/golymorph/objectpath"
)

// RulePolymorphism is a mapper that assigns a target type based on the given rules
type RulePolymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath

	// rules is a list of rules to apply. The first rule that matches is used to determine the target type.
	rules []Rule
}

func (p *RulePolymorphism) AssignTargetType(source any, target any) (error, bool) {

	// check for each rule if it matches and assign type if it does
	for _, rule := range p.rules {
		if err, matches := rule.Matches(source); err != nil {
			return errors.Join(errors.New("error applying rule"), err), false
		} else if matches {
			if err := objectpath.AssignTypeAtPath(target, p.targetPath, rule.NewType); err != nil {
				return errors.Join(errors.New("error assigning type to target"), err), false
			}
			return nil, true
		}
	}
	return nil, false

}
