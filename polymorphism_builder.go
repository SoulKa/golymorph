package golymorph

import "github.com/SoulKa/golymorph/objectpath"

type PolymorphismBuilder struct {
	targetPath objectpath.ObjectPath
	rules      []Rule
}

func NewPolymorphismBuilder() *PolymorphismBuilder {
	return &PolymorphismBuilder{*objectpath.NewSelfReferencePath(), []Rule{}}
}

func (b *PolymorphismBuilder) WithTargetPath(targetPath objectpath.ObjectPath) *PolymorphismBuilder {
	b.targetPath = targetPath
	return b
}
