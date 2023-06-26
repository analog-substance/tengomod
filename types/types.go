package types

import (
	"sort"

	"github.com/analog-substance/tengo/v2"
)

type Property struct {
	Get func() tengo.Object
	Set func(tengo.Object) error
}

func StaticProperty(val tengo.Object) Property {
	return Property{
		Get: func() tengo.Object {
			return val
		},
	}
}

type PropObject struct {
	tengo.ObjectImpl
	ObjectMap  map[string]tengo.Object
	Properties map[string]Property
}

// IndexGet gets an element at a given index
func (o *PropObject) IndexGet(index tengo.Object) (tengo.Object, error) {
	strIdx, ok := tengo.ToString(index)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	res, ok := o.ObjectMap[strIdx]
	if ok {
		return res, nil
	}
	res = tengo.UndefinedValue

	prop, ok := o.Properties[strIdx]
	if ok && prop.Get != nil {
		res = prop.Get()
	}
	return res, nil
}

// IndexSet sets an element at a given index.
func (o *PropObject) IndexSet(index tengo.Object, value tengo.Object) error {
	strIdx, ok := tengo.ToString(index)
	if !ok {
		return tengo.ErrInvalidIndexType
	}

	prop, ok := o.Properties[strIdx]
	if ok && prop.Set != nil {
		return prop.Set(value)
	}

	return nil
}

func (o *PropObject) Iterate() tengo.Iterator {
	var keys []string
	values := make(map[string]func() tengo.Object)
	for k := range o.ObjectMap {
		keys = append(keys, k)
		value := o.ObjectMap[k]
		values[k] = func() tengo.Object {
			return value
		}
	}

	for k, v := range o.Properties {
		keys = append(keys, k)
		values[k] = v.Get
	}

	sort.Strings(keys)

	return &PropObjectIterator{
		values: values,
		keys:   keys,
		length: len(keys),
	}
}

type PropObjectIterator struct {
	tengo.ObjectImpl
	values map[string]func() tengo.Object
	keys   []string
	i      int
	length int
}

// TypeName returns the name of the type.
func (i *PropObjectIterator) TypeName() string {
	return "prop-object-iterator"
}

func (i *PropObjectIterator) String() string {
	return "<prop-object-iterator>"
}

// IsFalsy returns true if the value of the type is falsy.
func (i *PropObjectIterator) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (i *PropObjectIterator) Equals(tengo.Object) bool {
	return false
}

// Copy returns a copy of the type.
func (i *PropObjectIterator) Copy() tengo.Object {
	return &PropObjectIterator{
		values: i.values,
		keys:   i.keys,
		i:      i.i,
		length: i.length,
	}
}

// Next returns true if there are more elements to iterate.
func (i *PropObjectIterator) Next() bool {
	i.i++
	return i.i <= i.length
}

// Key returns the key or index value of the current element.
func (i *PropObjectIterator) Key() tengo.Object {
	k := i.keys[i.i-1]
	return &tengo.String{Value: k}
}

// Value returns the value of the current element.
func (i *PropObjectIterator) Value() tengo.Object {
	k := i.keys[i.i-1]
	return i.values[k]()
}
