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

	tpl.Helpers.Add("say", func(name string) string {
		return fmt.Sprintf("say: %s", name)
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

func Test_Template_Clone(t *testing.T) {
	r := require.New(t)

	say := func(name string) string {
		return fmt.Sprintf("speak: %s", name)
	}

	input := `{{speak "mark"}}`
	t1, err := velvet.Parse(input)
	r.NoError(err)

	t2 := t1.Clone()
	t2.Helpers.Add("speak", say)

	ctx := velvet.NewContext()

	_, err = t1.Exec(ctx)
	r.Error(err)

	s, err := t2.Exec(ctx)
	r.NoError(err)
	r.Equal("speak: mark", s)
}
