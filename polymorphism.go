package golymorph

import (
	"errors"
	"fmt"
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
	"strings"
)

type Mapper interface {
}

// TypeMap is a map of discriminator values to reflect.Type
type TypeMap map[any]reflect.Type

// Polymorphism is a mapper that assigns a target type based on a discriminator value
type Polymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath

	// rules is a list of rules to apply. The first rule that matches is used to determine the target type.
	rules []Rule
}

type Polymorpher interface {
	AssignTargetType(source any, target any) error
}

// NewDiscriminatingPolymorphism creates a new Polymorphism mapper.
func NewDiscriminatingPolymorphism(discriminatorKey string, targetPath string, mapping TypeMap) (error, *Polymorphism) {

	// parse discriminator key
	err, discriminatorKeyObjectPath := objectpath.NewObjectPathFromString(discriminatorKey)
	if err != nil {
		return errors.Join(fmt.Errorf("error parsing discriminator key path"), err), nil
	}

	// parse target path
	if !strings.HasPrefix(targetPath, "/") {
		targetPath = "/" + targetPath // make target path absolute
	}
	err, targetObjectPath := objectpath.NewObjectPathFromString(targetPath)
	if err != nil {
		return errors.Join(fmt.Errorf("error parsing target path"), err), nil
	}

	// make discriminator path absolute
	if err := discriminatorKeyObjectPath.ToAbsolutePath(targetObjectPath); err != nil {
		return err, nil
	}

	// create rules
	var rules []Rule
	for discriminatorValue, targetType := range mapping {
		rules = append(rules, Rule{
			*discriminatorKeyObjectPath,
			ComparatorTypeEquality,
			discriminatorValue,
			targetType,
		})
	}

	return nil, &Polymorphism{
		*targetObjectPath,
		rules,
	}
}

// AssignTargetType assigns the determined type to target based on the polymorphism rules. The matching rule with the
// highest priority is used. If no rule matches, the target type is not changed. The source and target must be pointers.
func (polymorphism *Polymorphism) AssignTargetType(source any, target any) (error, bool) {

	// check for each rule if it matches and assign type if it does
	for _, rule := range polymorphism.rules {
		if err, matches := rule.Matches(source); err != nil {
			return errors.Join(errors.New("error applying rule"), err), false
		} else if matches {
			if err := objectpath.AssignTypeAtPath(target, polymorphism.targetPath, rule.NewType); err != nil {
				return errors.Join(errors.New("error assigning type to target"), err), false
			}
			return nil, true
		}
	}
	return nil, false

}
