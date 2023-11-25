package objectpath

import (
	"errors"
	"fmt"
	"testing"
)

func PathElementsFromStringArray(parts []string) []Element {
	var path []Element
	for _, part := range parts {
		path = append(path, Element{part, ElementTypeIdentifier})
	}
	return path
}

// TestParsePathString tests the ParsePathString function
func TestParsePathString(t *testing.T) {
	type TestCase struct {
		input  string
		output []Element
		error  error
	}
	testCases := []TestCase{
		{"/foo/bar", append([]Element{ElementRoot}, PathElementsFromStringArray([]string{"foo", "bar"})...), nil},
		{`foo/"bar"`, PathElementsFromStringArray([]string{"foo", "bar"}), nil},
		{`""/bar`, []Element{{"", ElementTypeIdentifier}, {"bar", ElementTypeIdentifier}}, nil},
		{"../bar/", []Element{ElementUpwardsReference, {"bar", ElementTypeIdentifier}}, nil},
		{"/../bar/..", []Element{ElementRoot, ElementUpwardsReference, {"bar", ElementTypeIdentifier}, ElementUpwardsReference}, nil},
		{`""`, []Element{{"", ElementTypeIdentifier}}, nil},
		{"./foo/./bar", []Element{ElementSelfReference, {"foo", ElementTypeIdentifier}, ElementSelfReference, {"bar", ElementTypeIdentifier}}, nil},
		{"foo/", []Element{{"foo", ElementTypeIdentifier}}, nil},
		{".", []Element{ElementSelfReference}, nil},
		{"", []Element{}, nil},
		{`foo/"bar`, nil, errors.New(`unexpected end of string after 8 runes. Expected """`)},
		{`fo#`, nil, errors.New(`unexpected character "#" at index 2. A non-enclosed path may only contain letters and digits`)},
		{`fo//bar`, nil, errors.New(`empty path element provided at index 3. Empty elements must be enclosed in quotes, e.g. /""/data`)},
		{`.../foo`, nil, errors.New(`invalid path element "..." at index 0. Only "." or ".." allowed`)},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`ParsePathString with input "%s"`, tc.input), func(t *testing.T) {
			var path Elements
			err := ParsePathString(tc.input, &path)
			if tc.error != nil {
				if err == nil || (tc.error.Error() != err.Error()) {
					t.Fatalf(`expected error "%s", but got "%s"`, tc.error, err)
				}
				return
			} else if err != nil {
				t.Fatalf("error parsing path: %s", err)
			}

			// compare output
			if tc.output == nil {
				t.Fatalf("invalid test. Expected either an error or a path, but got neither")
			} else if len(path) != len(tc.output) {
				t.Fatalf("expected length of path to be %d, but got %d", len(tc.output), len(path))
			}
			for i, part := range path {
				if part != tc.output[i] {
					t.Fatalf(`expected "%v" at path index %d, but got "%v"`, tc.output[i], i, part)
				}
			}
		})
	}
}

func TestNewObjectPathFromStringWithStringMethod(t *testing.T) {
	testCases := []string{`/"foo"/"bar"`, `"foo"/""/./..`}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf(`ParsePathString with input "%s"`, tc), func(t *testing.T) {
			err, path := NewObjectPathFromString(tc)
			if err != nil {
				t.Errorf("error parsing path: %s", err)
				return
			}
			for _, element := range path.elements {
				t.Log(element)
			}
			pathString := path.String()
			if pathString != tc {
				t.Errorf(`expected "%s", but got "%s"`, tc, pathString)
			}
		})
	}
}
