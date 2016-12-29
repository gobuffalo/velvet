package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_If_Helper(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `{{#if true}}hi{{/if}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("hi", s)
}

func Test_If_Helper_false(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `{{#if false}}hi{{/if}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("", s)
}

func Test_If_Helper_NoArgs(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `{{#if }}hi{{/if}}`

	_, err := velvet.Render(input, ctx)
	r.Error(err)
}

func Test_If_Helper_Else(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `
	{{#if false}}
		hi
	{{ else }}
		bye
	{{/if}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "bye")
}
