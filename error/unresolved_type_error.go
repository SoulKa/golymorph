package error

import "fmt"

// UnresolvedTypeError is an error that occurs when a type cannot be resolved when applying a polymorphism
type UnresolvedTypeError struct {
	Err        error
	TargetPath string
}

func (e *UnresolvedTypeError) Error() string {
	return fmt.Sprintf("unresolved type error at [%s]: %s", e.TargetPath, e.Err.Error())
}
