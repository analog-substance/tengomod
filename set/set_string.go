package set

import (
	"fmt"
	"sort"
	"strings"

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
		"add": &interop.AdvFunction{
			Name:    "add",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("item")},
			Value:   stringSet.add,
		},
		"add_range": &interop.AdvFunction{
			Name: "add_range",
			// NumArgs: interop.ExactArgs(1),
			Args:  []interop.AdvArg{interop.StrSliceArg("items", true)},
			Value: stringSet.addRange,
		},
		"slice": &tengo.UserFunction{
			Name:  "slice",
			Value: stringSet.slice,
		},
		"sorted_slice": &tengo.UserFunction{
			Name:  "sorted_slice",
			Value: stringSet.sortedSlice,
		},
	}

	stringSet.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return stringSet
}

func (s *StringSet) add(args interop.ArgMap) (tengo.Object, error) {
	item, _ := args.GetString("item")

	value := tengo.FalseValue
	if !s.Value.Contains(item) {
		s.Value.Add(item)
		value = tengo.TrueValue
	}

	return value, nil
}

func (s *StringSet) addRange(args interop.ArgMap) (tengo.Object, error) {
	items, _ := args.GetStringSlice("items")

	for _, item := range items {
		s.Value.Add(item)
	}

	return nil, nil
}

func (s *StringSet) slice(args ...tengo.Object) (tengo.Object, error) {
	var slice []string
	for _, val := range s.Value.Values() {
		slice = append(slice, val.(string))
	}

	return interop.GoStrSliceToTArray(slice), nil
}

func (s *StringSet) sortedSlice(args ...tengo.Object) (tengo.Object, error) {
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
	var elements []string
	for _, e := range s.Value.Values() {
		elements = append(elements, e.(string))
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
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
