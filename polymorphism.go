package golymorph

import (
	"errors"
	"fmt"
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

type TypeMap map[any]reflect.Type

type Polymorpher interface {
	// AssignTargetType assigns the determined type to target based on the polymorphism rules. The matching rule with the
	// highest priority is used. If no rule matches, the target type is not changed. The source and target must be pointers.
	AssignTargetType(source any, target any) (error, bool)
}

// RulePolymorphism is a mapper that assigns a target type based on the given rules
type RulePolymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath

	// rules is a list of rules to apply. The first rule that matches is used to determine the target type.
	rules []Rule
}

// TypeMapPolymorphism is a mapper that assigns a target type based on a discriminator value and a type map
type TypeMapPolymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath

	// discriminatorPath is the path to the discriminator value
	discriminatorPath objectpath.ObjectPath

	// typeMap is a map of discriminator values to types
	typeMap TypeMap
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

func (p *TypeMapPolymorphism) AssignTargetType(source any, target any) (error, bool) {

	// get discriminator value
	var discriminatorValue reflect.Value
	if err := objectpath.GetValueAtPath(source, p.discriminatorPath, &discriminatorValue); err != nil {
		return errors.Join(errors.New("error getting discriminator value"), err), false
	}
	rawDiscriminatorValue := discriminatorValue.Interface()
	fmt.Printf("discriminator value: %+v\n", rawDiscriminatorValue)

	// get type from type map
	if newType, ok := p.typeMap[rawDiscriminatorValue]; !ok {
		return nil, false
	} else if err := objectpath.AssignTypeAtPath(target, p.targetPath, newType); err != nil {
		return errors.Join(errors.New("error assigning type to target"), err), false
	}
	return nil, true
}
