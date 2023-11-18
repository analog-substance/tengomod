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
			Value:   stringSet.tengoAdd,
		},
		"add_range": &interop.AdvFunction{
			Name: "add_range",
			// NumArgs: interop.ExactArgs(1),
			Args:  []interop.AdvArg{interop.StrSliceArg("items", true)},
			Value: stringSet.tengoAddRange,
		},
		"slice": &tengo.UserFunction{
			Name:  "slice",
			Value: stringSet.tengoSlice,
		},
		"sorted_slice": &tengo.UserFunction{
			Name:  "sorted_slice",
			Value: stringSet.tengoSortedSlice,
		},
	}

	stringSet.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return stringSet
}

func (s *StringSet) tengoAdd(args interop.ArgMap) (tengo.Object, error) {
	item, _ := args.GetString("item")

	value := tengo.FalseValue
	if !s.add(item) {
		value = tengo.TrueValue
	}

	return value, nil
}

func (s *StringSet) add(item string) bool {
	if !s.Value.Contains(item) {
		s.Value.Add(item)
		return true
	}

	return false
}

func (s *StringSet) tengoAddRange(args interop.ArgMap) (tengo.Object, error) {
	items, _ := args.GetStringSlice("items")

	s.addRange(items)
	return nil, nil
}

func (s *StringSet) addRange(items []string) {
	for _, item := range items {
		s.Value.Add(item)
	}
}

func (s *StringSet) tengoSlice(args ...tengo.Object) (tengo.Object, error) {
	return interop.GoStrSliceToTArray(s.slice()), nil
}

func (s *StringSet) slice() []string {
	var slice []string
	for _, val := range s.Value.Values() {
		slice = append(slice, val.(string))
	}
	return slice
}

func (s *StringSet) tengoSortedSlice(args ...tengo.Object) (tengo.Object, error) {
	return interop.GoStrSliceToTArray(s.sortedSlice()), nil
}

func (s *StringSet) sortedSlice() []string {
	var slice []string
	for _, val := range s.Value.Values() {
		slice = append(slice, val.(string))
	}

	sort.Strings(slice)
	return slice
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
