package objectpath

// ElementType defines the type of Element
type ElementType int

const (
	// ElementTypeIdentifier is the type of Element that is a normal identifier
	ElementTypeIdentifier ElementType = iota
	// ElementTypeSelfReference is the type of Element that is a self reference ("."). It has no effect on the path.
	ElementTypeSelfReference
	// ElementTypeUpwardsReference is the type of Element that is an upwards reference ("..")
	ElementTypeUpwardsReference
	// ElementTypeRoot is the type of Element that is the root element
	ElementTypeRoot
)

// Element is a single element of a ObjectPath
type Element struct {
	name        string
	elementType ElementType
}

// MakeElement creates a new Element with the given name
func MakeElement(name string) Element {
	return Element{name, ElementTypeIdentifier}
}

// ElementSelfReference is a special Element that indicates a self reference. It has no effect on the path
var ElementSelfReference = Element{".", ElementTypeSelfReference}

// ElementUpwardsReference is a special Element that indicates an upwards reference
var ElementUpwardsReference = Element{"..", ElementTypeUpwardsReference}

// ElementRoot is a special Element that indicates the root element
var ElementRoot = Element{"", ElementTypeRoot}

// IsUpwardsReference returns true if the Element is the upward reference element
func (e *Element) IsUpwardsReference() bool {
	return e.elementType == ElementTypeUpwardsReference
}

// IsRootElement returns true if the Element is the root element
func (e *Element) IsRootElement() bool {
	return e.elementType == ElementTypeRoot
}
