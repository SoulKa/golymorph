package objectpath

import (
	"fmt"
	"unicode"
)

// ParsingState indicates the current state of the parser.
type ParsingState int

const (
	// The ParsingStateBeginning ParsingState is the state of the parser directly after a ParsingStateSlash or at the beginning of the string.
	ParsingStateBeginning ParsingState = iota
	// The ParsingStateName ParsingState is the state of the parser when parsing a non-enclosed name, e.g. `foo` in `foo/"bar"`.
	ParsingStateName
	// The ParsingStateEnclosedIdentifier ParsingState is the state of the parser when parsing an enclosed name, e.g. `"bar"` in `foo."bar"`.
	ParsingStateEnclosedIdentifier
	// The ParsingStateEscaping ParsingState is the state of the parser when parsing an escaped character, e.g. `"` after `\` in `foo/"bar\""`.
	ParsingStateEscaping
	// The ParsingStateSlash ParsingState is the state of the parser when parsing a slash, e.g. in `foo/bar`.
	ParsingStateSlash
	// The ParsingStateDot ParsingState is the state of the parser when parsing a dot, e.g. in `../bar`.
	ParsingStateDot
)

// Context contains the current ParsingState, input string, and other information needed for parsing.
type Context struct {
	chars     []rune
	path      *Elements
	state     ParsingState
	i         int
	reprocess bool
}

// hasChar returns true if the current index is within the bounds of the string.
func (ctx *Context) hasChar() bool {
	return ctx.i < len(ctx.chars)
}

// hasNextChar returns true if the next index is within the bounds of the string.
func (ctx *Context) hasNextChar() bool {
	return ctx.i+1 < len(ctx.chars)
}

// assertChar returns an error if the current character is not the expected character.
func (ctx *Context) assertChar(expected rune) error {
	if c := ctx.chars[ctx.i]; c != expected {
		return fmt.Errorf(`unexpected character [%c] at index %d. Expected [%c]`, c, ctx.i, expected)
	}
	return nil
}

// assertNextChar returns an error if the next character is not the expected character.
func (ctx *Context) assertNextChar(expected rune) error {
	if !ctx.hasNextChar() {
		return fmt.Errorf(`unexpected end of string after %d runes. Expected [%c]`, len(ctx.chars), expected)
	}
	ctx.i++
	return ctx.assertChar(expected)
}

// currentChar returns the current character.
func (ctx *Context) currentChar() rune {
	return ctx.chars[ctx.i]
}

// iterate increments the current index.
func (ctx *Context) iterate() {
	ctx.i++
}

// Parsing functions for each ParsingState.
var charParsingFunctions = []func(ctx *Context) error{

	ParsingStateBeginning: func(ctx *Context) error {
		ctx.path.appendElement()
		switch c := ctx.currentChar(); c {
		case '"':
			ctx.state = ParsingStateEnclosedIdentifier
		default:
			ctx.state = ParsingStateName
			ctx.reprocess = true // reprocess current character in ParsingStateName ParsingState
			return nil
		}
		return nil
	},

	ParsingStateName: func(ctx *Context) error {
		c := ctx.currentChar()
		switch c {
		case '/':
			if ctx.path.isCurrentPartEmpty() {
				return fmt.Errorf(`empty path element provided at index %d. Empty elements must be enclosed in quotes, e.g. /""/data`, ctx.i)
			}
			ctx.state = ParsingStateSlash
			ctx.reprocess = true
			return nil
		case '.':
			if ctx.path.isCurrentPartEmpty() {
				ctx.state = ParsingStateDot
				ctx.reprocess = true
				return nil
			}
		}
		if ctx.path.isCurrentPartEmpty() && !unicode.IsLetter(c) {
			return fmt.Errorf(`unexpected character [%c] at index %d. A non-enclosed path must start with a letter`, c, ctx.i)
		} else if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return fmt.Errorf(`unexpected character [%c] at index %d. A non-enclosed path may only contain letters and digits`, c, ctx.i)
		}
		ctx.path.appendCharToCurrentElement(c)
		return nil
	},

	ParsingStateEnclosedIdentifier: func(ctx *Context) error {
		switch c := ctx.currentChar(); c {
		case '\\':
			ctx.state = ParsingStateEscaping
		case '"':
			ctx.state = ParsingStateSlash
		default:
			ctx.path.appendCharToCurrentElement(c)
		}
		return nil
	},

	ParsingStateEscaping: func(ctx *Context) error {
		ctx.path.appendCharToCurrentElement(ctx.currentChar())
		ctx.state = ParsingStateEnclosedIdentifier
		return nil
	},

	ParsingStateSlash: func(ctx *Context) error {
		if err := ctx.assertChar('/'); err != nil {
			return err
		}
		ctx.state = ParsingStateBeginning
		return nil
	},

	ParsingStateDot: func(ctx *Context) error {
		e := ctx.path.currentElement()
		switch c := ctx.currentChar(); c {
		case '.':
			ctx.path.appendCharToCurrentElement(c)
			if e.name == "..." {
				return fmt.Errorf(`invalid path element [...] at index %d. Only [.] or [..] allowed`, ctx.i-2)
			} else if !ctx.hasNextChar() {
				ctx.chars = append(ctx.chars, '/') // append slash to end of string to parse last element
			}
		case '/':
			switch e.name {
			case ".":
				e.elementType = ElementTypeSelfReference
			case "..":
				e.elementType = ElementTypeUpwardsReference
			}
			ctx.state = ParsingStateSlash
			ctx.reprocess = true
		default:
			return fmt.Errorf(`unexpected character [%c] at index %d. Expected either [.] or [/]`, c, ctx.i)
		}
		return nil
	},
}

// ParsePathString parses a string into a slice of Element.
func ParsePathString(s string, path *Elements) error {
	*path = make(Elements, 0) // reset path
	ctx := Context{[]rune(s), path, ParsingStateBeginning, 0, false}

	// check if path is an absolute path
	if !ctx.hasChar() {
		return nil
	} else if ctx.currentChar() == '/' {
		ctx.path.appendElement()
		ctx.path.setCurrentElementType(ElementTypeRoot)
		ctx.iterate()
	}

	// parse while there are characters left
	for ctx.hasChar() {
		var charParsingFunction = charParsingFunctions[ctx.state]
		if err := charParsingFunction(&ctx); err != nil {
			return err
		}

		// increment index if not reprocessing current character
		if ctx.reprocess {
			ctx.reprocess = false
		} else {
			ctx.iterate()
		}
	}

	// check if current element is completely parsed
	if ctx.state == ParsingStateEnclosedIdentifier {
		return ctx.assertNextChar('"') // misuse methode to throw error
	} else if ctx.state == ParsingStateDot {
		ctx.chars = append(ctx.chars, '/') // append slash to end of string to parse last element
		if err := charParsingFunctions[ParsingStateDot](&ctx); err != nil {
			return err
		}
	}
	return nil
}
