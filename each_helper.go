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
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct || val.Len() == 0 {
		s, err := help.ElseBlock()
		return template.HTML(s), err
	}
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			ctx := help.Context.New()
			ctx.Set("@first", i == 0)
			ctx.Set("@last", i == val.Len()-1)
			ctx.Set("@index", i)
			ctx.Set("@value", v)
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
			ctx.Set("@last", i == len(keys)-1)
			ctx.Set("@key", key)
			ctx.Set("@value", v)
			s, err := help.BlockWith(ctx)
			if err != nil {
				return "", errors.WithStack(err)
			}
			out.WriteString(s)
		}
	}
	return template.HTML(out.String()), nil
}
