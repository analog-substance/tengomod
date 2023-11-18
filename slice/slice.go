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
			Value:   tengoContainsString,
		},
		"icontains_string": &interop.AdvFunction{
			Name:    "icontains_string",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false), interop.StrArg("input")},
			Value:   tengoIContainsString,
		},
		"rand_item": &interop.AdvFunction{
			Name:    "rand_item",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.SliceArg("slice", false)},
			Value:   tengoRandItem,
		},
		"unique": &interop.AdvFunction{
			Name:    "unique",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("slice", false)},
			Value:   tengoUnique,
		},
	}
}

func sortStrings(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	sort.Strings(slice)

	return interop.GoStrSliceToTArray(slice), nil
}

func tengoRandItem(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetSlice("slice")

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	item := randItem(slice, r1)
	if item == nil {
		return nil, nil
	}

	return item.(tengo.Object), nil
}

func randItem(slice []interface{}, r1 *rand.Rand) interface{} {
	if len(slice) == 0 {
		return nil
	}

	i := r1.Intn(len(slice))
	return slice[i]
}

func tengoUnique(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")

	return interop.GoStrSliceToTArray(unique(slice)), nil
}

func unique(slice []string) []string {
	set := hashset.New()
	for _, item := range slice {
		set.Add(item)
	}

	var items []string
	for _, item := range set.Values() {
		items = append(items, item.(string))
	}
	sort.Strings(items)

	return items
}

func tengoContainsString(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	input, _ := args.GetString("input")

	if containsString(slice, input) {
		return tengo.TrueValue, nil
	}

	return tengo.FalseValue, nil
}

func containsString(slice []string, input string) bool {
	for _, item := range slice {
		if item == input {
			return true
		}
	}

	return false
}

func tengoIContainsString(args interop.ArgMap) (tengo.Object, error) {
	slice, _ := args.GetStringSlice("slice")
	input, _ := args.GetString("input")

	if iContainsString(slice, input) {
		return tengo.TrueValue, nil
	}

	return tengo.FalseValue, nil
}

func iContainsString(slice []string, input string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, input) {
			return true
		}
	}

	return false
}
