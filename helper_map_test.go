package velvet_test

import (
	"testing"

	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_HelperMap_Add(t *testing.T) {
	r := require.New(t)
	hm, err := velvet.NewHelperMap()
	r.NoError(err)
	err = hm.Add("foo", func(help velvet.HelperContext) (string, error) {
		return "", nil
	})
	r.NoError(err)
	r.NotNil(hm.Helpers()["foo"])
}

func Test_HelperMap_Add_Invalid_NoReturn(t *testing.T) {
	r := require.New(t)

	hm, err := velvet.NewHelperMap()
	r.NoError(err)

	err = hm.Add("foo", func(help velvet.HelperContext) {})
	r.Error(err)
	r.Contains(err.Error(), "must return at least one")
	r.Nil(hm.Helpers()["foo"])
}

func Test_HelperMap_Add_Invalid_ReturnTypes(t *testing.T) {
	r := require.New(t)

	hm, err := velvet.NewHelperMap()
	r.NoError(err)

	err = hm.Add("foo", func(help velvet.HelperContext) (string, string) {
		return "", ""
	})
	r.Error(err)
	r.Contains(err.Error(), "foo must return ([string|template.HTML], [error]), not (string, string)")
	r.Nil(hm.Helpers()["foo"])

	err = hm.Add("foo", func(help velvet.HelperContext) int { return 1 })
	r.Error(err)
	r.Contains(err.Error(), "foo must return ([string|template.HTML], [error]), not (int)")
	r.Nil(hm.Helpers()["foo"])
}
