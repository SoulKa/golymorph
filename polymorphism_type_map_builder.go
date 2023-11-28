package golymorph

import (
	"errors"
	"github.com/SoulKa/golymorph/objectpath"
)

type polymorphismTypeMapBuilder struct {
	polymorphismBuilderBase
	typeMap           TypeMap
	discriminatorPath objectpath.ObjectPath
}

func (b *polymorphismTypeMapBuilder) WithDiscriminatorAt(discriminatorKey string) polymorphismBuilderFinalizer {
	if err, path := objectpath.NewObjectPathFromString(discriminatorKey); err != nil {
		b.errors = append(b.errors, err)
	} else if err := path.ToAbsolutePath(&b.targetPath); err != nil {
		b.errors = append(b.errors, err)
	} else {
		b.discriminatorPath = *path
	}
	return b
}

func (b *polymorphismTypeMapBuilder) Build() (error, TypeResolver) {
	if len(b.errors) > 0 {
		return errors.Join(b.errors...), nil
	}
	return nil, &TypeMapPolymorphism{
		Polymorphism: Polymorphism{
			TargetPath: b.targetPath},
		DiscriminatorPath: b.discriminatorPath,
		TypeMap:           b.typeMap}
}
