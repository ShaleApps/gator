package gator

import (
	"errors"
	"reflect"

	"github.com/ShaleApps/gator/Godeps/_workspace/src/github.com/onsi/gomega/matchers"
)

const (
	regexEmail    = `^([a-z0-9_\.-]+)@([\da-z\.-]+)\.([a-z\.]{2,6})$`
	regexHexColor = `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	regexUrl      = `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	regexIP       = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regexNum      = `^[1-9]\d*(\.\d+)?$`
	regexAlpha    = `^[a-zA-Z]*$`
)

type Func func(name string, v interface{}) error

type matcher interface {
	Match(actual interface{}) (success bool, err error)
}

var Nonzero = func(name string, v interface{}) error {
	m := &matchers.BeZeroMatcher{}
	zero, _ := m.Match(v)
	if zero {
		return formatError(name)
	}
	return nil
}

func Matches(regex string) Func {
	m := &matchers.MatchRegexpMatcher{Regexp: regex}
	return match(m)
}

var Email = Matches(regexEmail)
var HexColor = Matches(regexHexColor)
var URL = Matches(regexUrl)
var IP = Matches(regexIP)
var Alpha = Matches(regexAlpha)
var Num = Matches(regexNum)
var AlphaNum = combineFuncs(Matches("[a-zA-Z]+"), Matches("[0-9]+"))

func Gt(v interface{}) Func {
	return numericalMatch(">", v)
}

func Gte(v interface{}) Func {
	return numericalMatch(">=", v)
}

func Lt(v interface{}) Func {
	return numericalMatch("<", v)
}

func Lte(v interface{}) Func {
	return numericalMatch("<=", v)
}

var Lat = combineFuncs(
	numericalMatch("<=", 90.0),
	numericalMatch(">=", -90.0))

var Lon = combineFuncs(
	numericalMatch("<=", 180.0),
	numericalMatch(">=", -180.0))

func In(list []interface{}) Func {
	return func(k string, v interface{}) error {
		in := false
		for _, e := range list {
			in = in || reflect.DeepEqual(e, v)
		}
		if !in {
			return formatError(k)
		}
		return nil
	}
}

func NotIn(list []interface{}) Func {
	return func(k string, v interface{}) error {
		for _, e := range list {
			if reflect.DeepEqual(e, v) {
				return formatError(k)
			}
		}
		return nil
	}
}

func Len(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length != l {
			return formatError(k)
		}
		return nil
	}
}

func MinLen(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length < l {
			return formatError(k)
		}
		return nil
	}
}

func MaxLen(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length > l {
			return formatError(k)
		}
		return nil
	}
}

func Each(funcs ...Func) Func {
	return func(k string, v interface{}) error {
		if !isArrayOrSlice(v) {
			return formatError(k)
		}
		value := reflect.ValueOf(v)
		for i := 0; i < value.Len(); i++ {
			iFace := value.Index(i).Interface()
			for _, f := range funcs {
				if err := f(k, iFace); err != nil {
					return formatError(k)
				}
			}
		}
		return nil
	}
}

func numericalMatch(comparator string, v interface{}) Func {
	m := &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  []interface{}{v},
	}
	return match(m)
}

func match(m matcher) Func {
	return func(name string, v interface{}) error {
		matches, _ := m.Match(v)
		if !matches {
			return formatError(name)
		}
		return nil
	}
}

func combineFuncs(funcs ...Func) Func {
	return func(name string, v interface{}) error {
		for _, f := range funcs {
			if err := f(name, v); err != nil {
				return err
			}
		}
		return nil
	}
}

func formatError(name string) error {
	return errors.New(name + " did not pass validation.")
}
