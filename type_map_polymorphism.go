package golymorph

import (
	"errors"
	"fmt"
	golimorphError "github.com/SoulKa/golymorph/error"
	"github.com/SoulKa/golymorph/objectpath"
	"reflect"
)

// TypeMapPolymorphism is a mapper that assigns a target type based on a discriminator value and a type map
type TypeMapPolymorphism struct {
	Polymorphism

	// DiscriminatorPath is the path to the discriminator value
	DiscriminatorPath objectpath.ObjectPath

	// TypeMap is a map of discriminator values to types
	TypeMap TypeMap
}

func (p *TypeMapPolymorphism) AssignTargetType(source any, target any) error {

	// get discriminator value
	var discriminatorValue reflect.Value
	if err := objectpath.GetValueAtPath(source, p.DiscriminatorPath, &discriminatorValue); err != nil {
		return errors.Join(errors.New("error getting discriminator value"), err)
	}
	rawDiscriminatorValue := discriminatorValue.Interface()
	fmt.Printf("discriminator value: %+v\n", rawDiscriminatorValue)

	// get type from type map
	if newType, ok := p.TypeMap[rawDiscriminatorValue]; !ok {
		return &golimorphError.UnresolvedTypeError{
			Err:        fmt.Errorf("type map does not contain any key of value [%+v]", rawDiscriminatorValue),
			TargetPath: p.TargetPath.String(),
		}
	} else if err := objectpath.AssignTypeAtPath(target, p.TargetPath, newType); err != nil {
		return errors.Join(errors.New("error assigning type to target"), err)
	}
	return nil
}
