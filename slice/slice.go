package slice

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/emirpasic/gods/sets/hashset"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"sort_strings": &interop.AdvFunction{
			Name:    "sort_strings",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false)},
			Value:   sortStrings,
		},
		"contains_string": &interop.AdvFunction{
			Name:    "contains_string",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false), interop.StrArg("input")},
			Value:   containsString,
		},
		"icontains_string": &interop.AdvFunction{
			Name:    "icontains_string",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false), interop.StrArg("input")},
			Value:   iContainsString,
		},
		"rand_item": &interop.AdvFunction{
			Name:    "rand_item",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.SliceArg("slice", false)},
			Value:   randItem,
		},
		"unique": &interop.AdvFunction{
			Name:    "unique",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false)},
			Value:   unique,
		},
	}
}

func sortStrings(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	sort.Strings(slice)

	return interop.GoStrSliceToTArray(slice), nil
}

func randItem(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetSlice("slice")

	if len(slice) == 0 {
		return nil, nil
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	i := r1.Intn(len(slice))

	return slice[i].(tengo.Object), nil
}

func unique(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")

	set := hashset.New()
	for _, item := range slice {
		set.Add(item)
	}

	var items []string
	for _, item := range set.Values() {
		items = append(items, item.(string))
	}
	sort.Strings(items)

	return interop.GoStrSliceToTArray(items), nil
}

func containsString(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	input, _ := args.GetString("input")

	for _, item := range slice {
		if item == input {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}

func iContainsString(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	input, _ := args.GetString("input")

	for _, item := range slice {
		if strings.EqualFold(item, input) {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}
