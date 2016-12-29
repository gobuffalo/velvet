package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_Unless_Helper(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	input := `{{#unless false}}hi{{/unless}}`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("hi", s)
}
