package set

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/emirpasic/gods/sets/hashset"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"string_set": &interop.AdvFunction{
			Name:  "string_set",
			Args:  []interop.AdvArg{interop.StrSliceArg("items", true)},
			Value: newStringSet,
		},
	}
}

func newStringSet(args interop.ArgMap) (tengo.Object, error) {
	set := hashset.New()

	items, _ := args.GetStringSlice("items")
	for _, item := range items {
		set.Add(item)
	}

	return makeStringSet(set), nil
}
