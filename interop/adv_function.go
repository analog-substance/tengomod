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

	StrMapStrType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return TMapToGoStrMapStr(obj, name)
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

	ObjectType TypeValidator = func(obj tengo.Object, name string) (interface{}, error) {
		return obj, nil
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

func StrMapStrArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: StrMapStrType,
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

func ObjectArg(name string) AdvArg {
	return AdvArg{
		Name: name,
		Type: ObjectType,
	}
}

func UnionArg(name string, types ...TypeValidator) AdvArg {
	return AdvArg{
		Name: name,
		Type: UnionType(types...),
	}
}

func CustomArg(name string, t interface{}) AdvArg {
	return AdvArg{
		Name: name,
		Type: CustomType(t),
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

type ArgMap map[string]interface{}

func (m ArgMap) Exists(name string) bool {
	_, ok := m[name]
	return ok
}

func (m ArgMap) GetString(name string) (string, bool) {
	val, ok := m[name]
	if !ok {
		return "", ok
	}

	conv, ok := val.(string)
	return conv, ok
}

func (m ArgMap) GetStringSlice(name string) ([]string, bool) {
	val, ok := m[name]
	if !ok {
		return []string{}, ok
	}

	conv, ok := val.([]string)
	return conv, ok
}

func (m ArgMap) GetBool(name string) (bool, bool) {
	val, ok := m[name]
	if !ok {
		return false, ok
	}

	conv, ok := val.(bool)
	return conv, ok
}

func (m ArgMap) GetInt(name string) (int, bool) {
	val, ok := m[name]
	if !ok {
		return 0, ok
	}

	conv, ok := val.(int)
	return conv, ok
}

func (m ArgMap) GetIntSlice(name string) ([]int, bool) {
	val, ok := m[name]
	if !ok {
		return []int{}, ok
	}

	conv, ok := val.([]int)
	return conv, ok
}

func (m ArgMap) GetRegex(name string) (*regexp.Regexp, bool) {
	val, ok := m[name]
	if !ok {
		return nil, ok
	}

	conv, ok := val.(*regexp.Regexp)
	return conv, ok
}

func (m ArgMap) GetURL(name string) (*url.URL, bool) {
	val, ok := m[name]
	if !ok {
		return nil, ok
	}

	conv, ok := val.(*url.URL)
	return conv, ok
}

func (m ArgMap) GetCompiledFunc(name string) (*tengo.CompiledFunction, bool) {
	val, ok := m[name]
	if !ok {
		return nil, ok
	}

	conv, ok := val.(*tengo.CompiledFunction)
	return conv, ok
}

func (m ArgMap) GetSlice(name string) ([]interface{}, bool) {
	val, ok := m[name]
	if !ok {
		return []interface{}{}, ok
	}

	conv, ok := val.([]interface{})
	return conv, ok
}

func (m ArgMap) GetStrMapStr(name string) (map[string]string, bool) {
	val, ok := m[name]
	if !ok {
		return make(map[string]string), ok
	}

	conv, ok := val.(map[string]string)
	return conv, ok
}

func (m ArgMap) GetObject(name string) (tengo.Object, bool) {
	val, ok := m[name]
	if !ok {
		return nil, ok
	}

	conv, ok := val.(tengo.Object)
	return conv, ok
}

func (m ArgMap) Get(name string) (interface{}, bool) {
	val, ok := m[name]
	return val, ok
}

type AdvFunction struct {
	tengo.ObjectImpl
	Name    string
	NumArgs ArgValidator
	Args    []AdvArg
	Value   func(args ArgMap) (tengo.Object, error)
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

	args := make(ArgMap)
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
