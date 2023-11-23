package set_test

import (
	"sort"
	"testing"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/require"
	"github.com/analog-substance/tengomod/internal/test"
	"github.com/analog-substance/tengomod/set"
)

func toStringSlice(items []interface{}) []string {
	var slice []string
	for _, item := range items {
		slice = append(slice, item.(string))
	}

	return slice
}

func TestStringSet(t *testing.T) {
	stringSet := test.Module(t, "set").Call("string_set", "1", "2", "2", "3").Obj.(*set.StringSet)

	slice := toStringSlice(stringSet.Value.Values())
	sort.Strings(slice)

	require.Equal(t, []string{"1", "2", "3"}, slice)
}

func TestStringSetAdd(t *testing.T) {
	callRes := test.Module(t, "set").Call("string_set", "1")
	callRes.Call("add", "1").Expect(false)

	stringSet := callRes.Obj.(*set.StringSet)
	require.Equal(t, []string{"1"}, toStringSlice(stringSet.Value.Values()))

	callRes.Call("add", "2").Expect(true)

	slice := toStringSlice(stringSet.Value.Values())
	sort.Strings(slice)

	require.Equal(t, []string{"1", "2"}, slice)
}

func TestStringSetAddRange(t *testing.T) {
	callRes := test.Module(t, "set").Call("string_set", "1")
	callRes.Call("add_range", "1", "1", "2").ExpectNil()

	stringSet := callRes.Obj.(*set.StringSet)
	slice := toStringSlice(stringSet.Value.Values())
	sort.Strings(slice)

	require.Equal(t, []string{"1", "2"}, slice)
}

func TestStringSetSlice(t *testing.T) {
	callRes := test.Module(t, "set").Call("string_set", "1", "2", "3")

	obj := callRes.Call("slice").Obj
	require.IsType(t, &tengo.Array{}, obj)
	arr := obj.(*tengo.Array)

	var slice []string
	for _, item := range arr.Value {
		require.IsType(t, &tengo.String{}, item)
		slice = append(slice, item.(*tengo.String).Value)
	}

	sort.Strings(slice)
	require.Equal(t, []string{"1", "2", "3"}, slice)

	callRes.Call("sorted_slice").Expect(test.ARR{"1", "2", "3"})
}
