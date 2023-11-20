package slice_test

import (
	"testing"

	"github.com/analog-substance/tengomod/internal/test"
)

func TestSlice(t *testing.T) {
	test.Module(t, "slice").Call("sort_strings", test.ARR{
		"foo",
		"bar",
		"analog",
		"substance",
	}).Expect(test.ARR{
		"analog",
		"bar",
		"foo",
		"substance",
	})

	test.Module(t, "slice").Call("unique", test.ARR{
		"foo",
		"analog",
		"bar",
		"analog",
		"foo",
		"analog",
		"substance",
	}).Expect(test.ARR{
		"analog",
		"bar",
		"foo",
		"substance",
	})

	test.Module(t, "slice").Call("contains_string", test.ARR{
		"analog",
		"bar",
		"foo",
		"substance",
	}, "analog").Expect(true)
	test.Module(t, "slice").Call("contains_string", test.ARR{
		"analog",
		"bar",
		"foo",
		"substance",
	}, "example").Expect(false)

	test.Module(t, "slice").Call("icontains_string", test.ARR{
		"analog",
		"bar",
		"foo",
		"substance",
	}, "anaLoG").Expect(true)
}
