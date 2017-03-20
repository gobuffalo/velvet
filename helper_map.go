package velvet

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// HelperMap holds onto helpers and validates they are properly formed.
type HelperMap struct {
	moot    *sync.Mutex
	helpers map[string]interface{}
}

// NewHelperMap containing all of the "default" helpers from "velvet.Helpers".
func NewHelperMap() (HelperMap, error) {
	hm := HelperMap{
		helpers: map[string]interface{}{},
		moot:    &sync.Mutex{},
	}

	err := hm.AddMany(Helpers.Helpers())
	if err != nil {
		return hm, errors.WithStack(err)
	}
	return hm, nil
}

// Add a new helper to the map. New Helpers will be validated to ensure they
// meet the requirements for a helper:
/*
	func(...) (string) {}
	func(...) (string, error) {}
	func(...) (template.HTML) {}
	func(...) (template.HTML, error) {}
*/
func (h *HelperMap) Add(key string, helper interface{}) error {
	h.moot.Lock()
	defer h.moot.Unlock()
	if h.helpers == nil {
		h.helpers = map[string]interface{}{}
	}
	err := h.validateHelper(key, helper)
	if err != nil {
		return errors.WithStack(err)
	}
	h.helpers[key] = helper
	return nil
}

// AddMany helpers at the same time.
func (h *HelperMap) AddMany(helpers map[string]interface{}) error {
	for k, v := range helpers {
		err := h.Add(k, v)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Helpers returns the underlying list of helpers from the map
func (h HelperMap) Helpers() map[string]interface{} {
	return h.helpers
}

func (h *HelperMap) validateHelper(key string, helper interface{}) error {
	ht := reflect.ValueOf(helper).Type()

	if ht.NumOut() < 1 {
		return errors.WithStack(errors.Errorf("%s must return at least one value ([string|template.HTML], [error])", key))
	}
	so := ht.Out(0).Kind().String()
	if ht.NumOut() > 1 {
		et := ht.Out(1)
		ev := reflect.ValueOf(et)
		ek := fmt.Sprintf("%s", ev.Interface())
		if (so != "string" && so != "template.HTML") || (ek != "error") {
			return errors.WithStack(errors.Errorf("%s must return ([string|template.HTML], [error]), not (%s, %s)", key, so, et.Kind()))
		}
	} else {
		if so != "string" && so != "template.HTML" {
			return errors.WithStack(errors.Errorf("%s must return ([string|template.HTML], [error]), not (%s)", key, so))
		}
	}
	return nil
}
