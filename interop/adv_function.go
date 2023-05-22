package interop

import (
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/analog-substance/tengo/v2"
)

var (
	StrType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TStrToGoStr(obj, name)
	}

	IntType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TIntToGoInt(obj, name)
	}

	BoolType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TBoolToGoBool(obj, name)
	}

	StrSliceType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TArrayToGoStrSlice(obj, name)
	}

	IntSliceType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TArrayToGoIntSlice(obj, name)
	}

	SliceType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TArrayToGoInterfaceSlice(obj, name)
	}

	RegexType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		value, err := TStrToGoStr(obj, name)
		if err != nil {
			return nil, err
		}

		re, err := regexp.Compile(value)
		if err != nil {
			return GoErrToTErr(err), nil
		}

		return re, nil
	}

	URLType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		value, err := TStrToGoStr(obj, name)
		if err != nil {
			return nil, err
		}

		u, err := url.Parse(value)
		if err != nil {
			return GoErrToTErr(err), nil
		}

		return u, nil
	}

	CompileFuncType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		fn, ok := obj.(*tengo.CompiledFunction)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     name,
				Expected: "compiled-function",
				Found:    obj.TypeName(),
			}
		}
		return fn, nil
	}
)

type TypeValidator func(obj tengo.Object, name string) (interface{}, error)

func CustomType(t interface{}) TypeValidator {
	return func(obj tengo.Object, name string) (interface{}, error) {
		expectedType := reflect.TypeOf(t)
		v := reflect.ValueOf(obj)
		if !v.CanConvert(expectedType) {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     name,
				Expected: expectedType.Name(),
				Found:    obj.TypeName(),
			}
		}
		return v.Convert(expectedType).Interface(), nil
	}
}

func UnionType(types ...TypeValidator) TypeValidator {
	return func(obj tengo.Object, name string) (interface{}, error) {
		var possibleTypes []string
		for _, t := range types {
			value, err := t(obj, name)
			if err == nil {
				return value, nil
			}

			invalidType, ok := err.(tengo.ErrInvalidArgumentType)
			if ok {
				possibleTypes = append(possibleTypes, invalidType.Expected)
			} else {
				possibleTypes = append(possibleTypes, err.Error())
			}
		}
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: strings.Join(possibleTypes, "|"),
			Found:    obj.TypeName(),
		}
	}
}

func StrArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: StrType,
	}
}

func IntArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: IntType,
	}
}

func BoolArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: BoolType,
	}
}

func StrSliceArg(name string, varArgs bool) AdvArg {
	return AdvArg{
		Name:    name,
		Type:    StrSliceType,
		VarArgs: varArgs,
	}
}

func IntSliceArg(name string, varArgs bool) AdvArg {
	return AdvArg{
		Name:    name,
		Type:    IntSliceType,
		VarArgs: varArgs,
	}
}

func SliceArg(name string, varArgs bool) AdvArg {
	return AdvArg{
		Name:    name,
		Type:    SliceType,
		VarArgs: varArgs,
	}
}

func RegexArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: RegexType,
	}
}

func URLArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: URLType,
	}
}

func CompileFuncArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: CompileFuncType,
	}
}

func UnionArg(name string, types ...TypeValidator) AdvArg {
	return AdvArg{
		Name: name,
		Type: UnionType(types...),
	}
}

type AdvArg struct {
	Name    string
	Type    TypeValidator
	VarArgs bool
}

type ArgValidator func([]tengo.Object) error

func ExactArgs(n int) ArgValidator {
	return func(args []tengo.Object) error {
		if len(args) != n {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func MinArgs(min int) ArgValidator {
	return func(args []tengo.Object) error {
		if len(args) < min {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func MaxArgs(max int) ArgValidator {
	return func(args []tengo.Object) error {
		if len(args) > max {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func ArgRange(min int, max int) ArgValidator {
	return func(args []tengo.Object) error {
		if len(args) < min || len(args) > max {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

type AdvFunction struct {
	tengo.ObjectImpl
	Name    string
	NumArgs ArgValidator
	Args    []AdvArg
	Value   func(args map[string]interface{}) (tengo.Object, error)
}

// TypeName returns the name of the type.
func (o *AdvFunction) TypeName() string {
	return "adv-function:" + o.Name
}

func (o *AdvFunction) String() string {
	return "<adv-function>"
}

// Copy returns a copy of the type.
func (o *AdvFunction) Copy() tengo.Object {
	return &AdvFunction{
		Value:   o.Value,
		Name:    o.Name,
		NumArgs: o.NumArgs,
		Args:    o.Args,
	}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *AdvFunction) Equals(_ tengo.Object) bool {
	return false
}

// Call invokes a user function.
func (o *AdvFunction) Call(objs ...tengo.Object) (tengo.Object, error) {
	if o.NumArgs != nil {
		err := o.NumArgs(objs)
		if err != nil {
			return nil, err
		}
	}

	args := make(map[string]interface{})
	for i, arg := range o.Args {
		if i >= len(objs) {
			break
		}

		argObj := objs[i]
		if arg.VarArgs {
			argObj = &tengo.Array{
				Value: objs[i:],
			}
		}

		value, err := arg.Type(argObj, arg.Name)
		if err != nil {
			return nil, err
		}

		if errObj, ok := value.(*tengo.Error); ok {
			return errObj, nil
		}

		args[arg.Name] = value

		if arg.VarArgs {
			break
		}
	}
	return o.Value(args)
}

// CanCall returns whether the Object can be Called.
func (o *AdvFunction) CanCall() bool {
	return true
}
