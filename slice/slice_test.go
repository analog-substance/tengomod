package slice

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Test_sortStrings(t *testing.T) {
	type args struct {
		args interop.ArgMap
	}
	tests := []struct {
		name string
		args args
		want tengo.Object
	}{
		{
			name: "Sorted strings",
			args: args{
				args: interop.ArgMap{
					"slice": []string{
						"foo",
						"bar",
						"analog",
						"substance",
					},
				},
			},
			want: &tengo.Array{
				Value: []tengo.Object{
					&tengo.String{Value: "analog"},
					&tengo.String{Value: "bar"},
					&tengo.String{Value: "foo"},
					&tengo.String{Value: "substance"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := sortStrings(tt.args.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randItem(t *testing.T) {
	type args struct {
		slice []interface{}
		r1    *rand.Rand
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Empty slice",
			args: args{
				slice: []interface{}{},
				r1:    nil,
			},
			want: nil,
		},
		{
			name: "Random item",
			args: args{
				slice: []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				r1:    rand.New(rand.NewSource(2)),
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randItem(tt.args.slice, tt.args.r1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("randItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unique(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "With duplicates",
			args: args{
				slice: []string{
					"foo",
					"analog",
					"bar",
					"analog",
					"foo",
					"analog",
					"substance",
				},
			},
			want: []string{
				"analog",
				"bar",
				"foo",
				"substance",
			},
		},
		{
			name: "Without duplicates but not sorted",
			args: args{
				slice: []string{
					"foo",
					"bar",
					"analog",
					"substance",
				},
			},
			want: []string{
				"analog",
				"bar",
				"foo",
				"substance",
			},
		},
		{
			name: "Empty input",
			args: args{
				slice: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unique(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containsString(t *testing.T) {
	type args struct {
		slice []string
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Does contain string",
			args: args{
				slice: []string{
					"analog",
					"bar",
					"foo",
					"substance",
				},
				input: "analog",
			},
			want: true,
		},
		{
			name: "Doesn't contain string",
			args: args{
				slice: []string{
					"analog",
					"bar",
					"foo",
					"substance",
				},
				input: "example",
			},
			want: false,
		},
		{
			name: "Does contain string but different casing",
			args: args{
				slice: []string{
					"analog",
					"bar",
					"foo",
					"substance",
				},
				input: "Analog",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsString(tt.args.slice, tt.args.input); got != tt.want {
				t.Errorf("containsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iContainsString(t *testing.T) {
	type args struct {
		slice []string
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Does contain string",
			args: args{
				slice: []string{
					"analog",
					"bar",
					"foo",
					"substance",
				},
				input: "analog",
			},
			want: true,
		},
		{
			name: "Doesn't contain string",
			args: args{
				slice: []string{
					"analog",
					"bar",
					"foo",
					"substance",
				},
				input: "example",
			},
			want: false,
		},
		{
			name: "Does contain string but different casing",
			args: args{
				slice: []string{
					"anAloG",
					"bar",
					"foo",
					"substance",
				},
				input: "Analog",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iContainsString(tt.args.slice, tt.args.input); got != tt.want {
				t.Errorf("iContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
