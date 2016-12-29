package velvet

import (
	"bytes"
	"html/template"
	"reflect"

	"github.com/pkg/errors"
)

func eachHelper(collection interface{}, help HelperContext) (template.HTML, error) {
	out := bytes.Buffer{}
	val := reflect.ValueOf(collection)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			ctx := help.Context.New()
			ctx.Set("@first", i == 0)
			ctx.Set("@last", i == val.Len())
			ctx.Set("@index", i)
			ctx.Set("this", v)
			s, err := help.BlockWith(ctx)
			if err != nil {
				return "", errors.WithStack(err)
			}
			out.WriteString(s)
		}
	case reflect.Map:
		keys := val.MapKeys()
		for i := 0; i < len(keys); i++ {
			key := keys[i].Interface()
			v := val.MapIndex(keys[i]).Interface()
			ctx := help.Context.New()
			ctx.Set("@first", i == 0)
			ctx.Set("@last", i == len(keys))
			ctx.Set("@key", key)
			ctx.Set("this", v)
			s, err := help.BlockWith(ctx)
			if err != nil {
				return "", errors.WithStack(err)
			}
			out.WriteString(s)
		}
	}
	return template.HTML(out.String()), nil
}
