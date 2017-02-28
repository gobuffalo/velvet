package velvet

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"

	"github.com/aymerick/raymond/ast"
	"github.com/pkg/errors"
)

// HTMLer generates HTML source
type HTMLer interface {
	HTML() template.HTML
}

type interfacer interface {
	Interface() interface{}
}

var helperContextKind = "HelperContext"

type evalVisitor struct {
	template    *Template
	context     *Context
	curBlock    *ast.BlockStatement
	blockParams *blockParams
}

func newEvalVisitor(t *Template, c *Context) *evalVisitor {
	return &evalVisitor{
		template:    t,
		context:     c,
		blockParams: newBlockParams(),
	}
}

func (ev *evalVisitor) VisitProgram(p *ast.Program) interface{} {
	// fmt.Println("VisitProgram")
	defer ev.blockParams.pop()
	out := &bytes.Buffer{}
	ev.blockParams.push(p.BlockParams)
	for _, b := range p.Body {
		ev.context = ev.context.New()
		var value interface{}
		value = b.Accept(ev)
		switch vp := value.(type) {
		case error:
			return vp
		case template.HTML:
			out.Write([]byte(vp))
		case HTMLer:
			out.Write([]byte(vp.HTML()))
		case string:
			out.WriteString(template.HTMLEscapeString(vp))
		case []string:
			out.WriteString(template.HTMLEscapeString(strings.Join(vp, " ")))
		case int:
			out.WriteString(strconv.Itoa(vp))
		case fmt.Stringer:
			out.WriteString(template.HTMLEscapeString(vp.String()))
		case interfacer:
			out.WriteString(template.HTMLEscaper(vp.Interface()))
		case nil:
		default:
			return errors.WithStack(errors.Errorf("unsupport eval return format %T: %+v", value, value))
		}

	}
	return out.String()
}
func (ev *evalVisitor) VisitMustache(m *ast.MustacheStatement) interface{} {
	// fmt.Println("VisitMustache")
	expr := m.Expression.Accept(ev)
	return expr
}
func (ev *evalVisitor) VisitBlock(node *ast.BlockStatement) interface{} {
	// fmt.Println("VisitBlock")
	defer func() {
		ev.curBlock = nil
	}()
	ev.curBlock = node
	expr := node.Expression.Accept(ev)
	return expr
}

func (ev *evalVisitor) VisitPartial(*ast.PartialStatement) interface{} {
	// fmt.Println("VisitPartial")
	return ""
}

func (ev *evalVisitor) VisitContent(c *ast.ContentStatement) interface{} {
	// fmt.Println("VisitContent")
	return template.HTML(c.Original)
}

func (ev *evalVisitor) VisitComment(*ast.CommentStatement) interface{} {
	return ""
}

func (ev *evalVisitor) VisitExpression(e *ast.Expression) interface{} {
	// fmt.Println("VisitExpression")
	if e.Hash != nil {
		e.Hash.Accept(ev)
	}
	h := ev.helperName(e.HelperName())
	if h != "" {
		if helper, ok := ev.template.Helpers.Helpers()[h]; ok {
			return ev.evalHelper(e, helper)
		}
		if ev.context.Has(h) {
			x := ev.context.Get(h)
			if x != nil && h == "partial" {
				return ev.evalHelper(e, x)
			}
			return x
		}
		return errors.WithStack(errors.Errorf("could not find value for %s [line %d:%d]", h, e.Line, e.Pos))
	}
	parts := strings.Split(e.Canonical(), ".")
	if len(parts) > 1 && ev.context.Has(parts[0]) {
		rv := reflect.ValueOf(ev.context.Get(parts[0]))
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		m := rv.MethodByName(parts[1])
		if m.IsValid() {
			return ev.evalHelper(e, m.Interface())
		}
	}
	if fp := e.FieldPath(); fp != nil {
		return ev.VisitPath(fp)
	}
	if e.Path != nil {
		return e.Path.Accept(ev)
	}
	return nil
}

func (ev *evalVisitor) VisitSubExpression(*ast.SubExpression) interface{} {
	// fmt.Println("VisitSubExpression")
	return nil
}

func (ev *evalVisitor) VisitPath(node *ast.PathExpression) interface{} {
	// fmt.Println("VisitPath")
	// fmt.Printf("### node -> %+v\n", node)
	// fmt.Printf("### node -> %T\n", node)
	// fmt.Printf("### node.IsDataRoot() -> %+v\n", node.IsDataRoot())
	// fmt.Printf("### node.Loc() -> %+v\n", node.Location())
	// fmt.Printf("### node.String() -> %+v\n", node.String())
	// fmt.Printf("### node.Type() -> %+v\n", node.Type())
	// fmt.Printf("### node.Data -> %+v\n", node.Data)
	// fmt.Printf("### node.Depth -> %+v\n", node.Depth)
	// fmt.Printf("### node.Original -> %+v\n", node.Original)
	// fmt.Printf("### node.Parts -> %+v\n", node.Parts)
	// fmt.Printf("### node.Scoped -> %+v\n", node.Scoped)
	var v interface{}
	var h string
	if node.Data || len(node.Parts) == 0 {
		h = ev.helperName(node.Original)
	} else {
		h = ev.helperName(node.Parts[0])
	}
	if ev.context.Get(h) != nil {
		v = ev.context.Get(h)
	}
	if v == nil {
		return ""
		// 	return errors.WithStack(errors.Errorf("could not find value for %s [line %d:%d]", h, node.Line, node.Pos))
	}

	for i := 1; i < len(node.Parts); i++ {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		p := node.Parts[i]
		m := rv.MethodByName(p)
		if m.IsValid() {

			args := []reflect.Value{}
			rt := m.Type()
			if rt.NumIn() > 0 {
				last := rt.In(rt.NumIn() - 1)
				if last.Name() == helperContextKind {
					hargs := HelperContext{
						Context:     ev.context,
						Args:        []interface{}{},
						evalVisitor: ev,
					}
					args = append(args, reflect.ValueOf(hargs))
				} else if last.Kind() == reflect.Map {
					args = append(args, reflect.ValueOf(ev.context.Options()))
				}
				if len(args) > rt.NumIn() {
					err := errors.Errorf("Incorrect number of arguments being passed to %s (%d for %d)", p, len(args), rt.NumIn())
					return errors.WithStack(err)
				}
			}
			vv := m.Call(args)

			if len(vv) >= 1 {
				v = vv[0].Interface()
			}
			continue
		}
		switch rv.Kind() {
		case reflect.Map:
			pv := reflect.ValueOf(p)
			keys := rv.MapKeys()
			for i := 0; i < len(keys); i++ {
				k := keys[i]
				if k.Interface() == pv.Interface() {
					return rv.MapIndex(k).Interface()
				}
			}
			return errors.WithStack(errors.Errorf("could not find value for %s [line %d:%d]", node.Original, node.Line, node.Pos))
		default:
			f := rv.FieldByName(p)
			v = f.Interface()
		}
	}
	return v
}

func (ev *evalVisitor) VisitString(node *ast.StringLiteral) interface{} {
	// fmt.Println("VisitString")
	return node.Value
}

func (ev *evalVisitor) VisitBoolean(node *ast.BooleanLiteral) interface{} {
	// fmt.Println("VisitBoolean")
	return node.Value
}

func (ev *evalVisitor) VisitNumber(node *ast.NumberLiteral) interface{} {
	// fmt.Println("VisitNumber")
	return node.Number()
}

func (ev *evalVisitor) VisitHash(node *ast.Hash) interface{} {
	// fmt.Println("VisitHash")
	ctx := ev.context.New()
	for _, h := range node.Pairs {
		val := h.Accept(ev).(map[string]interface{})
		for k, v := range val {
			ctx.Set(k, v)
			ctx.Options()[k] = v
		}
	}
	ev.context = ctx
	return nil
}

func (ev *evalVisitor) VisitHashPair(node *ast.HashPair) interface{} {
	// fmt.Println("VisitHashPair")
	return map[string]interface{}{
		node.Key: node.Val.Accept(ev),
	}
}

func (ev *evalVisitor) evalHelper(node *ast.Expression, helper interface{}) (ret interface{}) {
	// fmt.Println("evalHelper")
	defer func() {
		if r := recover(); r != nil {
			switch rp := r.(type) {
			case error:
				ret = errors.WithStack(rp)
			case string:
				ret = errors.WithStack(errors.New(rp))
			}
		}
	}()

	hargs := HelperContext{
		Context:     ev.context,
		Args:        []interface{}{},
		evalVisitor: ev,
	}

	rv := reflect.ValueOf(helper)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()

	args := []reflect.Value{}

	if rt.NumIn() > 0 {
		for _, p := range node.Params {
			v := p.Accept(ev)
			vv := reflect.ValueOf(v)
			hargs.Args = append(hargs.Args, v)
			args = append(args, vv)
		}

		last := rt.In(rt.NumIn() - 1)
		if last.Name() == helperContextKind {
			args = append(args, reflect.ValueOf(hargs))
		} else if last.Kind() == reflect.Map {
			args = append(args, reflect.ValueOf(ev.context.Options()))
		}
		if len(args) > rt.NumIn() {
			err := errors.Errorf("Incorrect number of arguments being passed to %s (%d for %d)", node.Canonical(), len(args), rt.NumIn())
			return errors.WithStack(err)
		}
	}
	vv := rv.Call(args)

	if len(vv) >= 1 {
		v := vv[0].Interface()
		if len(vv) >= 2 {
			if !vv[1].IsNil() {
				return errors.WithStack(vv[1].Interface().(error))
			}
		}
		return v
	}

	return ""
}

func (ev *evalVisitor) helperName(h string) string {
	if h != "" {
		bp := ev.blockParams.current
		if len(bp) == 1 {
			if t := ev.context.Get("@value"); t != nil {
				ev.context.Set(bp[0], t)
			}
		}
		if len(bp) >= 2 {
			if t := ev.context.Get("@value"); t != nil {
				ev.context.Set(bp[1], t)
			}
			for _, k := range []string{"@index", "@key"} {
				if t := ev.context.Get(k); t != nil {
					ev.context.Set(bp[0], t)
				}
			}
		}
		return h
	}
	return ""
}
