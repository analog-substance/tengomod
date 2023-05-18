package types

import "github.com/analog-substance/tengo/v2"

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
