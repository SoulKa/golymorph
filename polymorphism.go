package golymorph

import (
	"github.com/SoulKa/golymorph/objectpath"
)

// Polymorphism is the base struct for all polymorphism mappers. It contains the target path to assign the new type to.
type Polymorphism struct {
	// TargetPath is the path to the object to assign the new type to
	TargetPath objectpath.ObjectPath
}
