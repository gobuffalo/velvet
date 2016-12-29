package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_MarkdownHelper(t *testing.T) {
	r := require.New(t)
	input := `{{markdown m}}`
	ctx := velvet.NewContext()
	ctx.Set("m", "# H1")
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "H1</h1>")
}
