package slice

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"sort_strings": &tengo.UserFunction{
			Name:  "sort_strings",
			Value: interop.NewCallable(sortStrings, interop.WithExactArgs(1)),
		},
		"contains_string": &tengo.UserFunction{
			Name:  "contains_string",
			Value: interop.NewCallable(containsString, interop.WithExactArgs(2)),
		},
		"icontains_string": &tengo.UserFunction{
			Name:  "icontains_string",
			Value: interop.NewCallable(iContainsString, interop.WithExactArgs(2)),
		},
		"rand_item": &tengo.UserFunction{
			Name:  "rand_item",
			Value: interop.NewCallable(randItem, interop.WithExactArgs(1)),
		},
		// "unique":          &tengo.UserFunction{Name: "unique", Value: interop.NewCallable(unique, interop.WithExactArgs(2))},
	}
}

func sortStrings(args ...tengo.Object) (tengo.Object, error) {
	slice, err := interop.TArrayToGoStringSlice(args[0], "slice")
	if err != nil {
		return nil, err
	}

	sort.Strings(slice)

	return interop.GoStringSliceToTArray(slice), nil
}

func randItem(args ...tengo.Object) (tengo.Object, error) {
	slice, err := interop.TArrayToGoInterfaceSlice(args[0], "slice")
	if err != nil {
		return nil, err
	}

	if len(slice) == 0 {
		return nil, nil
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	i := r1.Intn(len(slice))

	return slice[i].(tengo.Object), nil
}

// func unique(args ...tengo.Object) (tengo.Object, error) {
// 	if len(args) != 1 {
// 		return nil, tengo.ErrWrongNumArguments
// 	}

// 	array, ok := args[0].(*tengo.Array)
// 	if !ok {
// 		return nil, tengo.ErrInvalidArgumentType{
// 			Name:     "slice",
// 			Expected: "array",
// 			Found:    args[0].TypeName(),
// 		}
// 	}

// 	slice, err := arrayToStringSlice(array)
// 	if err != nil {
// 		return nil, err
// 	}

// 	itemSet := set.NewStringSet(slice)
// 	return sliceToStringArray(itemSet.SortedStringSlice()), nil
// }

func containsString(args ...tengo.Object) (tengo.Object, error) {
	slice, err := interop.TArrayToGoStringSlice(args[0], "slice")
	if err != nil {
		return nil, err
	}

	input, err := interop.TStringToGoString(args[1], "input")
	if err != nil {
		return nil, err
	}

	for _, item := range slice {
		if item == input {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}

func iContainsString(args ...tengo.Object) (tengo.Object, error) {
	slice, err := interop.TArrayToGoStringSlice(args[0], "slice")
	if err != nil {
		return nil, err
	}

	input, err := interop.TStringToGoString(args[1], "input")
	if err != nil {
		return nil, err
	}

	for _, item := range slice {
		if strings.EqualFold(item, input) {
			return tengo.TrueValue, nil
		}
	}
	return tengo.FalseValue, nil
}
