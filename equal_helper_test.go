package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_EqualHelper_True(t *testing.T) {
	r := require.New(t)
	input := `
	{{#eq 1 1}}
		it was true
	{{else}}
		it was false
	{{/eq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was true")
}

func Test_EqualHelper_False(t *testing.T) {
	r := require.New(t)
	input := `
	{{#eq 1 2}}
		it was true
	{{else}}
		it was false
	{{/eq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was false")
}

func Test_EqualHelper_DifferentTypes(t *testing.T) {
	r := require.New(t)
	input := `
	{{#eq 1 "1"}}
		it was true
	{{else}}
		it was false
	{{/eq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was false")
}

func Test_NotEqualHelper_True(t *testing.T) {
	r := require.New(t)
	input := `
	{{#neq 1 1}}
		it was true
	{{else}}
		it was false
	{{/neq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was false")
}

func Test_NotEqualHelper_False(t *testing.T) {
	r := require.New(t)
	input := `
	{{#neq 1 2}}
		it was true
	{{else}}
		it was false
	{{/neq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was true")
}

func Test_NotEqualHelper_DifferentTypes(t *testing.T) {
	r := require.New(t)
	input := `
	{{#neq 1 "1"}}
		it was true
	{{else}}
		it was false
	{{/neq}}
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "it was true")
}
