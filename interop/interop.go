package interop

import (
	"fmt"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/types"
)

func GoTSliceToGoInterfaceSlice(items []tengo.Object) []interface{} {
	var slice []interface{}
	for _, v := range items {
		slice = append(slice, v)
	}
	return slice
}

// GoStrSliceToTArray converts a golang string slice into a tengo Array
func GoStrSliceToTArray(slice []string) tengo.Object {
	var values []tengo.Object
	for _, s := range slice {
		values = append(values, &tengo.String{Value: s})
	}

	return &tengo.Array{
		Value: values,
	}
}

// GoIntSliceToTArray converts a golang int slice into a tengo Array
func GoIntSliceToTArray(slice []int) tengo.Object {
	var values []tengo.Object
	for _, i := range slice {
		values = append(values, &tengo.Int{Value: int64(i)})
	}

	return &tengo.Array{
		Value: values,
	}
}

// GoStrMapStrToTMap converts a golang map[string]string into a tengo Map
func GoStrMapStrToTMap(item map[string]string) tengo.Object {
	values := make(map[string]tengo.Object)
	for key, value := range item {
		values[key] = &tengo.String{Value: value}
	}

	return &tengo.Map{
		Value: values,
	}
}

// GoStrMapStrToTMap converts a golang map[string]string into a tengo ImmutableMap
func GoStrMapStrToTImmutMap(item map[string]string) tengo.Object {
	values := make(map[string]tengo.Object)
	for key, value := range item {
		values[key] = &tengo.String{Value: value}
	}

	return &tengo.ImmutableMap{
		Value: values,
	}
}

// TArrayToGoStrSlice converts a tengo Array into a golang string slice
func TArrayToGoStrSlice(obj tengo.Object, name string) ([]string, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoStrSlice(array.Value, name)
	case *tengo.ImmutableArray:
		return GoTSliceToGoStrSlice(array.Value, name)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoInterfaceSlice converts a tengo Array into a golang string slice
func TArrayToGoInterfaceSlice(obj tengo.Object, name string) ([]interface{}, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	case *tengo.ImmutableArray:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoTSlice converts a tengo Array into a golang tengo.Object slice
func TArrayToGoTSlice(obj tengo.Object, name string) ([]tengo.Object, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return array.Value, nil
	case *tengo.ImmutableArray:
		return array.Value, nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoSlice converts a tengo Array into a golang interface slice
func TArrayToGoSlice(obj tengo.Object, name string) ([]interface{}, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	case *tengo.ImmutableArray:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoIntSlice converts a tengo Array into a golang int slice
func TArrayToGoIntSlice(obj tengo.Object, name string) ([]int, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoIntSlice(array.Value, name)
	case *tengo.ImmutableArray:
		return GoTSliceToGoIntSlice(array.Value, name)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// GoTSliceToGoStrSlice converts a slice of tengo Objects into a golang string slice
func GoTSliceToGoStrSlice(slice []tengo.Object, name string) ([]string, error) {
	var strSlice []string
	for idx, obj := range slice {
		item, ok := tengo.ToString(obj)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s[%d]", name, idx),
				Expected: "string(compatible)",
				Found:    obj.TypeName(),
			}
		}
		strSlice = append(strSlice, item)
	}
	return strSlice, nil
}

// GoTSliceToGoIntSlice converts a slice of tengo Objects into a golang int slice
func GoTSliceToGoIntSlice(slice []tengo.Object, name string) ([]int, error) {
	var intSlice []int
	for idx, obj := range slice {
		i, ok := tengo.ToInt(obj)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%ss[%d]", name, idx),
				Expected: "int(compatible)",
				Found:    obj.TypeName(),
			}
		}

		intSlice = append(intSlice, i)
	}
	return intSlice, nil
}

// TMapToGoStrMapStr converts a tengo object into a golang map[string]string
func TMapToGoStrMapStr(obj tengo.Object, name string) (map[string]string, error) {
	var objMap map[string]tengo.Object
	switch o := obj.(type) {
	case *tengo.Map:
		objMap = o.Value
	case *tengo.ImmutableMap:
		objMap = o.Value
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "map(compatible)",
			Found:    obj.TypeName(),
		}
	}

	m := make(map[string]string)
	for key, value := range objMap {
		str, ok := tengo.ToString(value)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s key %s", name, key),
				Expected: "string(compatible)",
				Found:    value.TypeName(),
			}
		}

		m[key] = str
	}

	return m, nil
}

// TStrToGoStr converts a tengo object into a golang string
func TStrToGoStr(arg tengo.Object, name string) (string, error) {
	str, ok := tengo.ToString(arg)
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "string(compatible)",
			Found:    arg.TypeName(),
		}
	}

	return str, nil
}

// TIntToGoInt converts a tengo object into a golang int
func TIntToGoInt(arg tengo.Object, name string) (int, error) {
	i, ok := tengo.ToInt(arg)
	if !ok {
		return 0, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "int(compatible)",
			Found:    arg.TypeName(),
		}
	}

	return i, nil
}

// TBoolToGoBool converts a tengo object into a golang bool
func TBoolToGoBool(arg tengo.Object, name string) (bool, error) {
	b, ok := tengo.ToBool(arg)
	if !ok {
		return false, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "bool(compatible)",
			Found:    arg.TypeName(),
		}
	}

	return b, nil
}

// GoBoolToTBool converts a golang bool to a tengo bool
func GoBoolToTBool(val bool) tengo.Object {
	if val {
		return tengo.TrueValue
	}

	return tengo.FalseValue
}

// GoErrToTErr converts a golang error into a tengo Error
func GoErrToTErr(err error) tengo.Object {
	return &tengo.Error{
		Value: &tengo.String{
			Value: err.Error(),
		},
	}
}

// GoStrToTWarning converts a golang string into tengomod Warning
func GoStrToTWarning(value string) tengo.Object {
	return &types.Warning{
		Value: &tengo.String{
			Value: value,
		},
	}
}

// GoStrToTStr converts a golang string to a tengo string
func GoStrToTStr(str string) tengo.Object {
	return &tengo.String{
		Value: str,
	}
}

// GoIntToTInt converts a golang int to a tengo int
func GoIntToTInt(i int) tengo.Object {
	return &tengo.Int{
		Value: int64(i),
	}
}

// FuncASSSSRSp transform a function of 'func(string, string, string, string) *string' signature
// into tengo CallableFunc type.
func FuncASSSSRSp(fn func(string, string, string, string) *string) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(4),
		Args: []AdvArg{
			StrArg("first"),
			StrArg("second"),
			StrArg("third"),
			StrArg("fourth"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			s2, _ := args.GetString("second")
			s3, _ := args.GetString("third")
			s4, _ := args.GetString("fourth")

			s := fn(s1, s2, s3, s4)
			if len(*s) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			return &tengo.String{Value: *s}, nil
		},
	}

	return advFunc.Call
}

// FuncASSSRSp transform a function of 'func(string, string, string) *string' signature
// into tengo CallableFunc type.
func FuncASSSRSp(fn func(string, string, string) *string) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(3),
		Args: []AdvArg{
			StrArg("first"),
			StrArg("second"),
			StrArg("third"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			s2, _ := args.GetString("second")
			s3, _ := args.GetString("third")

			s := fn(s1, s2, s3)
			if len(*s) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			return &tengo.String{Value: *s}, nil
		},
	}
	return advFunc.Call
}

// FuncASSSsSRSsp transform a function of 'func(string, string, []string, string) *[]string' signature
// into tengo CallableFunc type.
func FuncASSSsSRSsp(fn func(string, string, []string, string) *[]string) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(4),
		Args: []AdvArg{
			StrArg("first"),
			StrArg("second"),
			StrSliceArg("third", false),
			StrArg("fourth"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			s2, _ := args.GetString("second")
			ss1, _ := args.GetStringSlice("third")
			s4, _ := args.GetString("fourth")

			return GoStrSliceToTArray(*fn(s1, s2, ss1, s4)), nil
		},
	}
	return advFunc.Call
}

// FuncASSsSRSsp transform a function of 'func(string, []string, string) *[]string' signature
// into tengo CallableFunc type.
func FuncASSsSRSsp(fn func(string, []string, string) *[]string) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(3),
		Args: []AdvArg{
			StrArg("first"),
			StrSliceArg("second", false),
			StrArg("third"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			ss1, _ := args.GetStringSlice("second")
			s3, _ := args.GetString("third")

			return GoStrSliceToTArray(*fn(s1, ss1, s3)), nil
		},
	}
	return advFunc.Call
}

// FuncASSBSRBp transform a function of 'func(string, string, bool, string) *bool' signature
// into tengo CallableFunc type.
func FuncASSBSRBp(fn func(string, string, bool, string) *bool) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(4),
		Args: []AdvArg{
			StrArg("first"),
			StrArg("second"),
			BoolArg("third"),
			StrArg("fourth"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			s2, _ := args.GetString("second")
			b1, _ := args.GetBool("third")
			s4, _ := args.GetString("fourth")

			return GoBoolToTBool(*fn(s1, s2, b1, s4)), nil
		},
	}
	return advFunc.Call
}

// FuncASBSRBp transform a function of 'func(string, bool, string) *bool' signature
// into tengo CallableFunc type.
func FuncASBSRBp(fn func(string, bool, string) *bool) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(3),
		Args: []AdvArg{
			StrArg("first"),
			BoolArg("second"),
			StrArg("third"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			b1, _ := args.GetBool("second")
			s3, _ := args.GetString("third")

			return GoBoolToTBool(*fn(s1, b1, s3)), nil
		},
	}
	return advFunc.Call
}

// FuncASSISRIp transform a function of 'func(string, string, int, string) *int' signature
// into tengo CallableFunc type.
func FuncASSISRIp(fn func(string, string, int, string) *int) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(4),
		Args: []AdvArg{
			StrArg("first"),
			StrArg("second"),
			IntArg("third"),
			StrArg("fourth"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			s2, _ := args.GetString("second")
			i1, _ := args.GetInt("third")
			s4, _ := args.GetString("fourth")

			i := fn(s1, s2, i1, s4)
			return &tengo.Int{Value: int64(*i)}, nil
		},
	}
	return advFunc.Call
}

// FuncASISRIp transform a function of 'func(string, int, string) *int' signature
// into tengo CallableFunc type.
func FuncASISRIp(fn func(string, int, string) *int) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(3),
		Args: []AdvArg{
			StrArg("first"),
			IntArg("second"),
			StrArg("third"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			i1, _ := args.GetInt("second")
			s3, _ := args.GetString("third")

			i := fn(s1, i1, s3)
			return &tengo.Int{Value: int64(*i)}, nil
		},
	}
	return advFunc.Call
}

// FuncASRSsE transform a function of 'func(string) ([]string, error)' signature
// into tengo CallableFunc type.
func FuncASRSsE(fn func(string) ([]string, error)) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(1),
		Args: []AdvArg{
			StrArg("first"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")

			res, err := fn(s1)
			if err != nil {
				return GoErrToTErr(err), nil
			}

			return GoStrSliceToTArray(res), nil
		},
	}
	return advFunc.Call
}

// FuncASRBE transform a function of 'func(string) (bool, error)' signature
// into tengo CallableFunc type.
func FuncASRBE(fn func(string) (bool, error)) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(1),
		Args: []AdvArg{
			StrArg("first"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")

			res, err := fn(s1)
			if err != nil {
				return GoErrToTErr(err), nil
			}

			return GoBoolToTBool(res), nil
		},
	}
	return advFunc.Call
}

// FuncASRB transform a function of 'func(string) bool' signature
// into tengo CallableFunc type.
func FuncASRB(fn func(string) bool) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(1),
		Args: []AdvArg{
			StrArg("first"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			return GoBoolToTBool(fn(s1)), nil
		},
	}
	return advFunc.Call
}

// FuncASvRSsE transform a function of 'func(...string) ([]string, error)' signature
// into tengo CallableFunc type.
func FuncASvRSsE(fn func(...string) ([]string, error)) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: MinArgs(1),
		Args: []AdvArg{
			StrSliceArg("first", true),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			strings, _ := args.GetStringSlice("first")

			res, err := fn(strings...)
			if err != nil {
				return GoErrToTErr(err), nil
			}

			return GoStrSliceToTArray(res), nil
		},
	}
	return advFunc.Call
}

// FuncASvRB transform a function of 'func(...string) bool' signature
// into tengo CallableFunc type.
func FuncASvRB(fn func(...string) bool) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: MinArgs(1),
		Args: []AdvArg{
			StrSliceArg("first", true),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			strings, _ := args.GetStringSlice("first")

			return GoBoolToTBool(fn(strings...)), nil
		},
	}
	return advFunc.Call
}

// FuncASvRS transform a function of 'func(...string) string' signature
// into tengo CallableFunc type.
func FuncASvRS(fn func(...string) string) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: MinArgs(1),
		Args: []AdvArg{
			StrSliceArg("first", true),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			strings, _ := args.GetStringSlice("first")

			return GoStrToTStr(fn(strings...)), nil
		},
	}
	return advFunc.Call
}

// FuncASRI transform a function of 'func(string) int' signature into
// CallableFunc type.
func FuncASRI(fn func(string) int) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(1),
		Args: []AdvArg{
			StrArg("first"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			s1, _ := args.GetString("first")
			return GoIntToTInt(fn(s1)), nil
		},
	}
	return advFunc.Call
}

// FuncABR transform a function of 'func(bool)' signature into
// CallableFunc type.
func FuncABR(fn func(bool)) tengo.CallableFunc {
	advFunc := AdvFunction{
		NumArgs: ExactArgs(1),
		Args: []AdvArg{
			BoolArg("first"),
		},
		Value: func(args ArgMap) (tengo.Object, error) {
			b1, _ := args.GetBool("first")
			fn(b1)
			return nil, nil
		},
	}
	return advFunc.Call
}

// AliasFunc is used to call the same tengo function using a different name
func AliasFunc(obj tengo.Object, name string, src string) *tengo.UserFunction {
	return &tengo.UserFunction{
		Name: name,
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			fn, err := obj.IndexGet(&tengo.String{Value: src})
			if err != nil {
				return nil, err
			}
			return fn.Call(args...)
		},
	}
}
