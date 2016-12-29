package velvet_test

import (
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_CustomGlobalHelper(t *testing.T) {
	r := require.New(t)
	err := velvet.Helpers.Add("say", func(name string) (string, error) {
		return fmt.Sprintf("say: %s", name), nil
	})
	r.NoError(err)

	input := `{{say "mark"}}`
	ctx := velvet.NewContext()
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("say: mark", s)
}

func Test_CustomGlobalBlockHelper(t *testing.T) {
	r := require.New(t)
	velvet.Helpers.Add("say", func(name string, help velvet.HelperContext) (template.HTML, error) {
		ctx := help.Context
		ctx.Set("name", strings.ToUpper(name))
		s, err := help.BlockWith(ctx)
		return template.HTML(s), err
	})

	input := `
	{{#say "mark"}}
		<h1>{{name}}</h1>
	{{/say}}
	`
	ctx := velvet.NewContext()
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Contains(s, "<h1>MARK</h1>")
}

func Test_Helper_Hash_Options(t *testing.T) {
	r := require.New(t)
	velvet.Helpers.Add("say", func(help velvet.HelperContext) string {
		return help.Context.Get("name").(string)
	})

	input := `{{say name="mark"}}`
	ctx := velvet.NewContext()
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("mark", s)
}

func Test_Helper_Hash_Options_Many(t *testing.T) {
	r := require.New(t)
	velvet.Helpers.Add("say", func(help velvet.HelperContext) string {
		return help.Context.Get("first").(string) + help.Context.Get("last").(string)
	})

	input := `{{say first=first_name last=last_name}}`
	ctx := velvet.NewContext()
	ctx.Set("first_name", "Mark")
	ctx.Set("last_name", "Bates")
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal("MarkBates", s)
}

func Test_Helper_Santize_Output(t *testing.T) {
	r := require.New(t)

	velvet.Helpers.Add("safe", func(help velvet.HelperContext) template.HTML {
		return template.HTML("<p>safe</p>")
	})
	velvet.Helpers.Add("unsafe", func(help velvet.HelperContext) string {
		return "<b>unsafe</b>"
	})

	input := `{{safe}}|{{unsafe}}`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Equal("<p>safe</p>|&lt;b&gt;unsafe&lt;/b&gt;", s)
}

func Test_JSON_Helper(t *testing.T) {
	r := require.New(t)

	input := `{{json names}}`
	ctx := velvet.NewContext()
	ctx.Set("names", []string{"mark", "bates"})
	s, err := velvet.Render(input, ctx)
	r.NoError(err)
	r.Equal(`["mark","bates"]`, s)
}
