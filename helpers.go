package velvet

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Helpers contains all of the default helpers for velvet.
// These will be available to all templates. You should add
// any custom global helpers to this list.
var Helpers = HelperMap{}

func init() {
	Helpers.Add("if", ifHelper)
	Helpers.Add("unless", unlessHelper)
	Helpers.Add("each", eachHelper)
	Helpers.Add("json", toJSONHelper)
	Helpers.Add("js_escape", template.JSEscapeString)
	Helpers.Add("html_escape", template.HTMLEscapeString)
	Helpers.Add("upcase", strings.ToUpper)
	Helpers.Add("downcase", strings.ToLower)
	Helpers.Add("content_for", contentForHelper)
	Helpers.Add("content_of", contentOfHelper)
	Helpers.Add("markdown", markdownHelper)
	Helpers.Add("debug", debugHelper)
	Helpers.AddMany(inflect.Helpers)
}

// HelperContext is an optional context that can be passed
// as the last argument to helper functions.
type HelperContext struct {
	Context     *Context
	Args        []interface{}
	evalVisitor *evalVisitor
}

// Block executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement.
func (h HelperContext) Block() (string, error) {
	return h.BlockWith(h.Context)
}

// BlockWith executes the block of template associated with
// the helper, think the block inside of an "if" or "each"
// statement. It takes a new context with which to evaluate
// the block.
func (h HelperContext) BlockWith(ctx *Context) (string, error) {
	nev := newEvalVisitor(h.evalVisitor.template, ctx)
	nev.blockParams = h.evalVisitor.blockParams
	dd := nev.VisitProgram(h.evalVisitor.curBlock.Program)
	switch tp := dd.(type) {
	case string:
		return tp, nil
	case error:
		return "", errors.WithStack(tp)
	case nil:
		return "", nil
	default:
		return "", errors.WithStack(errors.Errorf("unknown return value %T %+v", dd, dd))
	}
}

// ElseBlock executes the "inverse" block of template associated with
// the helper, think the "else" block of an "if" or "each"
// statement.
func (h HelperContext) ElseBlock() (string, error) {
	return h.ElseBlockWith(h.Context)
}

// ElseBlockWith executes the "inverse" block of template associated with
// the helper, think the "else" block of an "if" or "each"
// statement. It takes a new context with which to evaluate
// the block.
func (h HelperContext) ElseBlockWith(ctx *Context) (string, error) {
	if h.evalVisitor.curBlock.Inverse == nil {
		return "", nil
	}
	nev := newEvalVisitor(h.evalVisitor.template, ctx)
	nev.blockParams = h.evalVisitor.blockParams
	dd := nev.VisitProgram(h.evalVisitor.curBlock.Inverse)
	switch tp := dd.(type) {
	case string:
		return tp, nil
	case error:
		return "", errors.WithStack(tp)
	case nil:
		return "", nil
	default:
		return "", errors.WithStack(errors.Errorf("unknown return value %T %+v", dd, dd))
	}
}

// toJSONHelper converts an interface into a string.
func toJSONHelper(v interface{}) (template.HTML, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return template.HTML(b), nil
}

// Debug by verbosely printing out using 'pre' tags.
func debugHelper(v interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("<pre>%+v</pre>", v))
}
