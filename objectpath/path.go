package objectpath

import (
	"errors"
	"fmt"
)

// Elements is a slice of Element
type Elements []Element

// appendElement appends a new, empty Element to the path
func (path *Elements) appendElement() {
	*path = append(*path, Element{"", ElementTypeIdentifier})
}

// currentElement returns a pointer to the last element of the path
func (path *Elements) currentElement() *Element {
	return &(*path)[len(*path)-1]
}

// appendCharToCurrentElement appends a character to the ParsingStateName of the last element of the path
func (path *Elements) appendCharToCurrentElement(char rune) {
	path.currentElement().name += string(char)
}

// isCurrentPartEmpty returns true if the last element of the path has an empty ParsingStateName
func (path *Elements) isCurrentPartEmpty() bool {
	return path.currentElement().name == ""
}

// setCurrentElementType sets the type of the last element of the path
func (path *Elements) setCurrentElementType(elementType ElementType) {
	path.currentElement().elementType = elementType
}

// ObjectPath manages Elements and may be absolute or relative
type ObjectPath struct {
	elements   Elements
	isAbsolute bool
}

// NewObjectPathFromString creates a new ObjectPath from a string
func NewObjectPathFromString(s string) (error, *ObjectPath) {
	var path ObjectPath
	if err := ParsePathString(s, &path.elements); err != nil {
		return err, nil
	}

	// check if path is absolute
	if path.getLength() > 0 && path.elements[0].IsRootElement() {
		path.isAbsolute = true
		if err := path.DeleteAt(0, 1); err != nil {
			return err, nil
		}
	}
	return nil, &path
}

// NewSelfReferencePath creates a new ObjectPath with a single self reference element. The path is relative.
func NewSelfReferencePath() *ObjectPath {
	return &ObjectPath{Elements{ElementSelfReference}, false}
}

// IsAbsolutePath returns true if the path starts with a root element
func (p *ObjectPath) IsAbsolutePath() bool {
	return p.isAbsolute
}

// IsRelativePath returns true if the path does not start with a root element
func (p *ObjectPath) IsRelativePath() bool {
	return !p.IsAbsolutePath()
}

// Normalize an absolute path by resolving upwards references
func (p *ObjectPath) Normalize() error {
	if !p.IsAbsolutePath() {
		return fmt.Errorf("cannot normalize a relative path")
	}
	for i := 0; i < p.getLength(); i++ {
		part := p.elements[i]
		switch part.elementType {
		case ElementTypeSelfReference:
			if err := p.DeleteAt(i, 1); err != nil {
				return err
			}
			i--
		case ElementTypeUpwardsReference:
			if p.getLength() == 0 {
				return errors.New("path escaping root")
			}
			if err := p.DeleteAt(i-1, 2); err != nil {
				return err
			}
			i -= 2
		}
	}
	return nil
}

// ToAbsolutePath converts a relative path to an absolute path using the given absolute reference path
func (p *ObjectPath) ToAbsolutePath(referencePath *ObjectPath) error {
	if p.IsAbsolutePath() {
		return nil
	}
	if !referencePath.IsAbsolutePath() {
		return errors.New("the given reference path must be absolute")
	}

	// concatenate the paths
	p.elements = append(referencePath.elements, p.elements...)
	p.isAbsolute = true
	if err := p.Normalize(); err != nil {
		return fmt.Errorf("error normalizing path: %s", err)
	}
	return nil
}

// Pop removes the last element of the path
func (p *ObjectPath) Pop() error {
	if p.getLength() == 0 {
		return fmt.Errorf("cannot pop from empty path")
	}
	p.elements = p.elements[:len(p.elements)-1]
	return nil
}

// Push appends an element to the path
func (p *ObjectPath) Push(element Element) error {
	p.elements = append(p.elements, element)
	return nil
}

// DeleteAt removes n elements starting at index
func (p *ObjectPath) DeleteAt(index int, n int) error {
	if p.getLength() < index+n || index < 0 {
		return fmt.Errorf("invalid index %d for path of length %d", index, p.getLength())
	}
	p.elements = append(p.elements[:index], p.elements[index+n:]...)
	return nil
}

// IsEqualTo returns true if the path is equal to the given path
func (p *ObjectPath) IsEqualTo(other *ObjectPath) bool {
	if p.getLength() != other.getLength() {
		return false
	}
	for i, part := range p.elements {
		if part != other.elements[i] {
			return false
		}
	}
	return true
}

// getLength returns the length of the path
func (p *ObjectPath) getLength() int {
	return len(p.elements)
}

// String returns the string representation of the path. This string would lead to the same path when parsed again.
func (p *ObjectPath) String() string {
	var s string
	if p.isAbsolute {
		s = "/"
	}
	for i, part := range p.elements {
		if i > 0 {
			s += "/"
		}
		if part.elementType == ElementTypeIdentifier {
			s += string('"')
			s += part.name
			s += string('"')
		} else {
			s += part.name
		}
	}
	return s
}
