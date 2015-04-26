package gator

import (
	"fmt"
	"net/url"
	"reflect"
)

type Validator interface {
	Validate() error
}

type Gator struct {
	vals []Validator
}

func New() *Gator {
	return &Gator{vals: []Validator{}}
}

func (g *Gator) Add(v Validator, vals ...Validator) *Gator {
	g.vals = append(g.vals, v)
	g.vals = append(g.vals, vals...)
	return g
}

func (g *Gator) Validate() error {
	for _, v := range g.vals {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Field struct {
	name string
	src  interface{}
	f    Func
}

func NewField(name string, src interface{}, f Func) *Field {
	return &Field{
		name: name,
		src:  src,
		f:    f,
	}
}

func (f *Field) Validate() error {
	return f.f(f.name, f.src)
}

type QueryStr struct {
	str string
	src interface{}
}

func NewQueryStr(src interface{}, str string) *QueryStr {
	return &QueryStr{src: src, str: str}
}

func (q *QueryStr) Validate() error {
	m, err := url.ParseQuery(q.str)
	if err != nil {
		fmt.Errorf("gator: couldn't parse QueryStr - %s", err)
	}

	if err := isStructOrStructPtr(q.src); err != nil {
		return err
	}
	objT := reflect.TypeOf(q.src)
	objV := reflect.ValueOf(q.src)

	g := New()
	for i := 0; i < objT.NumField(); i++ {
		name := objT.Field(i).Name
		value := objV.Field(i).Interface()

		for key, values := range m {
			if key == name {
				for _, v := range values {
					for _, f := range funcsFromTag(v) {
						g.Add(NewField(name, value, f))
					}
				}
			}
		}
	}
	return g.Validate()
}
