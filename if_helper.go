package velvet

import (
	"html/template"
	"reflect"
)

func ifHelper(conditional interface{}, help HelperContext) (template.HTML, error) {
	if IsTrue(conditional) {
		s, err := help.Block()
		return template.HTML(s), err
	}
	s, err := help.ElseBlock()
	return template.HTML(s), err
}

// IsTrue returns true if obj is a truthy value.
func IsTrue(obj interface{}) bool {
	thruth, ok := isTrueValue(reflect.ValueOf(obj))
	if !ok {
		return false
	}
	return thruth
}

// isTrueValue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value
//
// NOTE: borrowed from https://github.com/golang/go/tree/master/src/text/template/exec.go
func isTrueValue(val reflect.Value) (truth, ok bool) {
	if !val.IsValid() {
		// Something like var x interface{}, never set. It's a form of nil.
		return false, true
	}
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
		truth = !val.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Struct:
		truth = true // Struct values are always true.
	default:
		return
	}
	return truth, true
}
