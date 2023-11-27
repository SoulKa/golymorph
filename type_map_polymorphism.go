package golymorph

import (
	"errors"
	"fmt"
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

// TypeMapPolymorphism is a mapper that assigns a target type based on a discriminator value and a type map
type TypeMapPolymorphism struct {
	// targetPath is the path to the object to assign the new type to
	targetPath objectpath.ObjectPath

	// discriminatorPath is the path to the discriminator value
	discriminatorPath objectpath.ObjectPath

	// typeMap is a map of discriminator values to types
	typeMap TypeMap
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
