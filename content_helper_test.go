package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_ContentForOf(t *testing.T) {
	r := require.New(t)
	input := `
	{{#content_for "buttons"}}<button>hi</button>{{/content_for}}
	<b1>{{content_of "buttons"}}</b1>
	<b2>{{content_of "buttons"}}</b2>
	`
	s, err := velvet.Render(input, velvet.NewContext())
	r.NoError(err)
	r.Contains(s, "<b1><button>hi</button></b1>")
	r.Contains(s, "<b2><button>hi</button></b2>")
}
