package velvet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_blockParams(t *testing.T) {
	r := require.New(t)
	bp := newBlockParams()
	r.Equal([]string{}, bp.current)
	r.Len(bp.stack, 0)

	bp.push([]string{"mark"})
	r.Equal([]string{"mark"}, bp.current)
	r.Len(bp.stack, 1)

	bp.push([]string{"bates"})
	r.Equal([]string{"bates"}, bp.current)
	r.Len(bp.stack, 2)
	r.Equal([][]string{
		[]string{"mark"},
		[]string{"bates"},
	}, bp.stack)

	b := bp.pop()
	r.Equal([]string{"bates"}, b)
	r.Equal([]string{"mark"}, bp.current)
	r.Len(bp.stack, 1)

	b = bp.pop()
	r.Equal([]string{"mark"}, b)
	r.Len(bp.stack, 0)
	r.Equal([]string{}, bp.current)
}
