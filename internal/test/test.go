package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/require"
	"github.com/analog-substance/tengo/v2/stdlib"
	"github.com/analog-substance/tengomod"
)

type ARR = []interface{}
type MAP = map[string]interface{}
type IARR []interface{}
type IMAP map[string]interface{}

type CallRes struct {
	t   *testing.T
	Obj interface{}
	Err error
}

func (c CallRes) Call(funcName string, args ...interface{}) CallRes {
	if c.Err != nil {
		return c
	}

	var oargs []tengo.Object
	for _, v := range args {
		oargs = append(oargs, Object(v))
	}

	switch o := c.Obj.(type) {
	case *tengo.BuiltinModule:
		m, ok := o.Attrs[funcName]
		if !ok {
			return CallRes{t: c.t, Err: fmt.Errorf(
				"function not found: %s", funcName)}
		}

		var res tengo.Object
		var err error
		f, ok := m.(*tengo.UserFunction)
		if ok {
			res, err = f.Value(oargs...)
		} else if m.CanCall() {
			res, err = m.Call(oargs...)
		} else {
			return CallRes{t: c.t, Err: fmt.Errorf(
				"non-callable: %s", funcName)}
		}

		return CallRes{t: c.t, Obj: res, Err: err}
	case *tengo.UserFunction:
		res, err := o.Value(oargs...)
		return CallRes{t: c.t, Obj: res, Err: err}
	case *tengo.ImmutableMap:
		m, ok := o.Value[funcName]
		if !ok {
			return CallRes{t: c.t, Err: fmt.Errorf("function not found: %s", funcName)}
		}

		f, ok := m.(*tengo.UserFunction)
		if !ok {
			return CallRes{t: c.t, Err: fmt.Errorf("non-callable: %s", funcName)}
		}

		res, err := f.Value(oargs...)
		return CallRes{t: c.t, Obj: res, Err: err}
	default:
		if obj, ok := c.Obj.(tengo.Object); ok && obj.CanCall() {
			res, err := obj.Call(oargs...)
			return CallRes{t: c.t, Obj: res, Err: err}
		}

		panic(fmt.Errorf("unexpected object: %v (%T)", o, o))
	}
}

func (c CallRes) Expect(expected interface{}, msgAndArgs ...interface{}) {
	require.NoError(c.t, c.Err, msgAndArgs...)
	require.Equal(c.t, Object(expected), c.Obj, msgAndArgs...)
}

func (c CallRes) ExpectError() {
	require.Error(c.t, c.Err)
}

func (c CallRes) ExpectTengoError() {
	require.IsType(c.t, &tengo.Error{}, c.Obj)
}

func Module(t *testing.T, moduleName string) CallRes {
	mod := tengomod.GetModuleMap(tengomod.WithModules(moduleName)).GetBuiltinModule(moduleName)
	if mod == nil {
		return CallRes{t: t, Err: fmt.Errorf("module not found: %s", moduleName)}
	}

	return CallRes{t: t, Obj: mod}
}

func Object(v interface{}) tengo.Object {
	if v == nil {
		return nil
	}

	switch v := v.(type) {
	case tengo.Object:
		return v
	case string:
		return &tengo.String{Value: v}
	case int64:
		return &tengo.Int{Value: v}
	case int: // for convenience
		return &tengo.Int{Value: int64(v)}
	case bool:
		if v {
			return tengo.TrueValue
		}
		return tengo.FalseValue
	case rune:
		return &tengo.Char{Value: v}
	case byte: // for convenience
		return &tengo.Char{Value: rune(v)}
	case float64:
		return &tengo.Float{Value: v}
	case []byte:
		return &tengo.Bytes{Value: v}
	case MAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = Object(v)
		}

		return &tengo.Map{Value: objs}
	case ARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, Object(e))
		}

		return &tengo.Array{Value: objs}
	case IMAP:
		objs := make(map[string]tengo.Object)
		for k, v := range v {
			objs[k] = Object(v)
		}

		return &tengo.ImmutableMap{Value: objs}
	case IARR:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, Object(e))
		}

		return &tengo.ImmutableArray{Value: objs}
	case time.Time:
		return &tengo.Time{Value: v}
	case []int:
		var objs []tengo.Object
		for _, e := range v {
			objs = append(objs, &tengo.Int{Value: int64(e)})
		}

		return &tengo.Array{Value: objs}
	}

	panic(fmt.Errorf("unknown type: %T", v))
}

func Expect(t *testing.T, input string, expected interface{}) {
	s := tengo.NewScript([]byte(input))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	c, err := s.Run()
	require.NoError(t, err)
	require.NotNil(t, c)
	v := c.Get("out")
	require.NotNil(t, v)
	require.Equal(t, expected, v.Value())
}
