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
	discriminatorKey objectpath.ObjectPath
	targetPath       objectpath.ObjectPath
	mapping          TypeMap
}

// NewPolymorphism creates a new Polymorphism mapper
func NewPolymorphism(discriminatorKey string, mapping TypeMap) (error, *Polymorphism) {
	return NewPolymorphismAtPath(discriminatorKey, "/", mapping)
}

// NewPolymorphismAtPath creates a new Polymorphism mapper.
func NewPolymorphismAtPath(discriminatorKey string, targetPath string, mapping TypeMap) (error, *Polymorphism) {

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

	return nil, &Polymorphism{
		*discriminatorKeyObjectPath,
		*targetObjectPath,
		mapping,
	}
}

// AssignTargetType assigns the target type based on the discriminator value in the source.
// The source and target must be pointers.
func (polymorphism *Polymorphism) AssignTargetType(source any, target any) error {

	// get discriminator value
	var discriminatorVal reflect.Value
	if err := objectpath.GetValueAtPath(source, polymorphism.discriminatorKey, &discriminatorVal); err != nil {
		return errors.Join(errors.New("error getting discriminator value"), err)
	}
	discriminatorValue := discriminatorVal.Interface()

	// get type for discriminator value
	targetType, ok := polymorphism.mapping[discriminatorValue]
	if !ok {
		return fmt.Errorf("no target type found for discriminator value %s", discriminatorValue)
	}

	// create target with type
	if err := objectpath.AssignTypeAtPath(target, polymorphism.targetPath, targetType); err != nil {
		return errors.Join(errors.New("error assigning type to target"), err)
	}
	return nil

}
