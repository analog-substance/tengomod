package set

import (
	"sort"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
	"github.com/emirpasic/gods/sets/hashset"
)

type StringSet struct {
	types.PropObject
	Value *hashset.Set
}

func makeStringSet(set *hashset.Set) *StringSet {
	stringSet := &StringSet{
		Value: set,
	}

	objectMap := map[string]tengo.Object{
		"add": &tengo.UserFunction{
			Name:  "add",
			Value: interop.NewCallable(stringSet.add, interop.WithExactArgs(1)),
		},
		"add_range": &tengo.UserFunction{
			Name:  "add_range",
			Value: interop.NewCallable(stringSet.addRange, interop.WithExactArgs(1)),
		},
		"sorted_string_slice": &tengo.UserFunction{
			Name:  "sorted_string_slice",
			Value: stringSet.sortedStringSlice,
		},
	}

	stringSet.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return stringSet
}

func (s *StringSet) add(args ...tengo.Object) (tengo.Object, error) {
	item, err := interop.TStrToGoStr(args[0], "item")
	if err != nil {
		return nil, err
	}

	value := tengo.FalseValue
	if !s.Value.Contains(item) {
		s.Value.Add(item)
		value = tengo.TrueValue
	}

	return value, nil
}

func (s *StringSet) addRange(args ...tengo.Object) (tengo.Object, error) {
	items, err := interop.TArrayToGoStrSlice(args[0], "items")
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		s.Value.Add(item)
	}

	return nil, nil
}

func (s *StringSet) sortedStringSlice(args ...tengo.Object) (tengo.Object, error) {
	var slice []string
	for _, val := range s.Value.Values() {
		slice = append(slice, val.(string))
	}

	sort.Strings(slice)

	return interop.GoStrSliceToTArray(slice), nil
}

// TypeName should return the name of the type.
func (s *StringSet) TypeName() string {
	return "string-set"
}

// String should return a string representation of the type's value.
func (s *StringSet) String() string {
	return s.Value.String()
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (s *StringSet) IsFalsy() bool {
	return s.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (s *StringSet) CanIterate() bool {
	return false
}
