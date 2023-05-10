package set

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/emirpasic/gods/sets/hashset"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"new_string_set": &tengo.UserFunction{Name: "new_string_set", Value: newStringSet},
	}
}

func newStringSet(args ...tengo.Object) (tengo.Object, error) {
	set := hashset.New()
	return makeStringSet(set), nil
}
