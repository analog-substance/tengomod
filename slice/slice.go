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

func sortStrings(args map[string]interface{}) (tengo.Object, error) {
	slice := args["slice"].([]string)
	sort.Strings(slice)

	return interop.GoStrSliceToTArray(slice), nil
}

func randItem(args map[string]interface{}) (tengo.Object, error) {
	slice := args["slice"].([]interface{})

	if len(slice) == 0 {
		return nil, nil
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	i := r1.Intn(len(slice))

	return slice[i].(tengo.Object), nil
}

func unique(args map[string]interface{}) (tengo.Object, error) {
	slice := args["slice"].([]string)

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

func containsString(args map[string]interface{}) (tengo.Object, error) {
	slice := args["slice"].([]string)
	input := args["input"].(string)

	for _, item := range slice {
		if item == input {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}

func iContainsString(args map[string]interface{}) (tengo.Object, error) {
	slice := args["slice"].([]string)
	input := args["input"].(string)

	for _, item := range slice {
		if strings.EqualFold(item, input) {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}
