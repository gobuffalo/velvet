package velvet

import "github.com/pkg/errors"

var cache = map[string]*Template{}

// Parse an input string and return a Template.
func Parse(input string) (*Template, error) {
	if t, ok := cache[input]; ok {
		return t, nil
	}
	t, err := NewTemplate(input)

	if err == nil {
		cache[input] = t
	}

	if err != nil {
		return t, errors.WithStack(err)
	}

	return t, nil
}

// Render a string using the given the context.
func Render(input string, ctx *Context) (string, error) {
	t, err := Parse(input)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return t.Exec(ctx)
}
