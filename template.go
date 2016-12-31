package velvet

import (
	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/parser"
	"github.com/pkg/errors"
)

// Template represents an input and helpers to be used
// to evaluate and render the input.
type Template struct {
	Input   string
	Helpers HelperMap
	program *ast.Program
}

// NewTemplate from the input string. Adds all of the
// global helper functions from "velvet.Helpers".
func NewTemplate(input string) (*Template, error) {
	hm, err := NewHelperMap()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	t := &Template{
		Input:   input,
		Helpers: hm,
	}
	err = t.Parse()
	if err != nil {
		return t, errors.WithStack(err)
	}
	return t, nil
}

// Parse the template this can be called many times
// as a sucessful result is cached and is used on subsequent
// uses.
func (t *Template) Parse() error {
	if t.program != nil {
		return nil
	}
	program, err := parser.Parse(t.Input)
	if err != nil {
		return errors.WithStack(err)
	}
	t.program = program
	return nil
}

// Exec the template using the content and return the results
func (t *Template) Exec(ctx *Context) (string, error) {
	err := t.Parse()
	if err != nil {
		return "", errors.WithStack(err)
	}
	v := newEvalVisitor(t, ctx)
	r := t.program.Accept(v)
	switch rp := r.(type) {
	case string:
		return rp, nil
	case error:
		return "", rp
	case nil:
		return "", nil
	default:
		return "", errors.WithStack(errors.Errorf("unsupport eval return format %T: %+v", r, r))
	}
}

// Clone a template. This is useful for defining helpers on per "instance" of the template.
func (t *Template) Clone() *Template {
	hm, _ := NewHelperMap()
	hm.AddMany(t.Helpers.Helpers())
	t2 := &Template{
		Helpers: hm,
		Input:   t.Input,
		program: t.program,
	}
	return t2
}
