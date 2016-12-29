package velvet_test

import (
	"strings"
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_Eval_Map_Call_Key(t *testing.T) {
	r := require.New(t)
	ctx := velvet.NewContext()
	data := map[string]string{
		"a": "A",
		"b": "B",
	}
	ctx.Set("letters", data)
	input := `
	{{letters.a}}|{{letters.b}}
	`

	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("A|B", strings.TrimSpace(s))
}

func Test_Eval_Calls_on_Pointers(t *testing.T) {
	r := require.New(t)
	type user struct {
		Name string
	}
	u := &user{Name: "Mark"}
	ctx := velvet.NewContext()
	ctx.Set("user", u)

	s, err := velvet.Render("{{user.Name}}", ctx)
	r.NoError(err)
	r.Equal("Mark", s)
}
