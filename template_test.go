package velvet_test

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_Template_Helpers(t *testing.T) {
	r := require.New(t)

	input := `{{say "mark"}}`
	tpl, err := velvet.Parse(input)
	r.NoError(err)

	tpl.Helpers.Add("say", func(name string, help velvet.HelperContext) (string, error) {
		return fmt.Sprintf("say: %s", name), nil
	})

	ctx := velvet.NewContext()
	s, err := tpl.Exec(ctx)
	r.NoError(err)
	r.Equal("say: mark", s)

	input = `{{say "jane"}}`
	tpl, err = velvet.Parse(input)
	r.NoError(err)
	_, err = tpl.Exec(ctx)
	r.Error(err)
}
