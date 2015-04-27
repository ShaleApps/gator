package gator

import (
	"errors"
	"reflect"

	"github.com/ShaleApps/gator/Godeps/_workspace/src/github.com/onsi/gomega/matchers"
)

const (
	regexEmail    = `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`
	regexHexColor = `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	regexURL      = `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	regexIP       = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regexNum      = `^[1-9]\d*(\.\d+)?$`
	regexAlpha    = `^[a-zA-Z]*$`
)

// Func is a validation function that returns an error if v is invalid.
type Func func(name string, v interface{}) error

// Matches returns a Func that validates against the given regex.
func Matches(regex string) Func {
	m := &matchers.MatchRegexpMatcher{Regexp: regex}
	return match(m)
}

// Nonzero returns a Func that validates its value is non-zero.
// http://golang.org/pkg/reflect/#Zero
func Nonzero() Func {
	return func(name string, v interface{}) error {
		m := &matchers.BeZeroMatcher{}
		zero, _ := m.Match(v)
		if zero {
			return formatError(name)
		}
		return nil
	}
}

// Email returns a Func that validates its value is an email address.
func Email() Func {
	return Matches(regexEmail)
}

// HexColor returns a Func that validates its value is a hexidecimal number prefixed by a hash.
// HTML standard link: http://www.w3.org/TR/REC-html40/types.html#h-6.5
func HexColor() Func {
	return Matches(regexHexColor)
}

// URL returns a Func that validates its value is a URL.
func URL() Func {
	return Matches(regexURL)
}

// IP returns a Func that validates its value is an IP address.
func IP() Func {
	return Matches(regexIP)
}

// Alpha returns a Func that validates its value contains only letters.
func Alpha() Func {
	return Matches(regexAlpha)
}

// Num returns a Func that validates its value contains only numbers.
func Num() Func {
	return Matches(regexNum)
}

// AlphaNum returns a Func that validates its value contains both numbers and letters.
func AlphaNum() Func {
	return combineFuncs(Matches("[a-zA-Z]+"), Matches("[0-9]+"))
}

// Gt returns a Func that validates its value is a number greater than v.
func Gt(v interface{}) Func {
	return numericalMatch(">", v)
}

// Gte returns a Func that validates its value is a number greater than or equal to v.
func Gte(v interface{}) Func {
	return numericalMatch(">=", v)
}

// Lt returns a Func that validates its value is a number less than v.
func Lt(v interface{}) Func {
	return numericalMatch("<", v)
}

// Lte returns a Func that validates its value is a number less than or equal to v.
func Lte(v interface{}) Func {
	return numericalMatch("<=", v)
}

// Lat returns a Func that validates its value is a decimal between 90 and -90.
func Lat() Func {
	return combineFuncs(
		numericalMatch("<=", 90.0),
		numericalMatch(">=", -90.0))
}

// Lon returns a Func that validates its value is a decimal between 180 and -180.
func Lon() Func {
	return combineFuncs(
		numericalMatch("<=", 180.0),
		numericalMatch(">=", -180.0))
}

// In returns a Func that validates its value is in the inputed list.  Comparisons
// use reflect.DeepEqual.
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

// NotIn returns a Func that validates its value is not in the inputed list.  Comparisons
// use reflect.DeepEqual.
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

// Len returns a Func that validates its value's length is equal to l.
func Len(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length != l {
			return formatError(k)
		}
		return nil
	}
}

// MinLen returns a Func that validates its value's length is greater than or equal to l.
func MinLen(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length < l {
			return formatError(k)
		}
		return nil
	}
}

// MaxLen returns a Func that validates its value's length is less than or equal to l.
func MaxLen(l int) Func {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length > l {
			return formatError(k)
		}
		return nil
	}
}

// Each returns a Func that validates the list of functions for each element in an array or slice.
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

type matcher interface {
	Match(actual interface{}) (success bool, err error)
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
