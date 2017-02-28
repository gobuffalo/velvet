package velvet

import (
	"sync"

	"github.com/gobuffalo/buffalo/render"
	"github.com/pkg/errors"
)

func BuffaloRenderer(input string, options render.TemplateOptions) (string, error) {
	t, err := Parse(input)
	if err != nil {
		return "", err
	}
	t.Helpers.AddMany(options.Helpers)
	return t.Exec(NewContextWith(options.Data))
}

var cache = map[string]*Template{}
var moot = &sync.Mutex{}

// Parse an input string and return a Template.
func Parse(input string) (*Template, error) {
	moot.Lock()
	defer moot.Unlock()
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
