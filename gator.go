package gator

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	structTagKey = "gator"
)

// Validator is the interface that wraps the basic Validate method.
type Validator interface {
	Validate() error
}

// Gator is a Validator that is comprised of other Validators.
type Gator struct {
	vals []Validator
}

// New creates an initialized Gator.
func New(options ...func(*Gator)) *Gator {
	return &Gator{vals: []Validator{}}
}

// NewStruct generates validation fields based on src's gator struct
// tags and adds them to the returned gator.  If src isn't a struct
// or pointer to a struct an error will be returned from the Validation
// method.
func NewStruct(src interface{}) *Gator {
	g := New()
	objT, objV, err := getReflectInfo(src)
	if err != nil {
		g.Add(errValidator{err: err})
		return g
	}

	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		tag := field.Tag.Get(structTagKey)
		name := field.Name
		funcs := funcsFromTag(tag)
		value := objV.Field(i).Interface()

		for _, f := range funcs {
			g.Add(NewField(name, value, f))
		}
	}
	return g
}

// NewQueryStr generates validation fields by parsing queryStr using
// url.ParseQuery and adds them to the returned gator.  If the queryStr
// can't be parsed or if src isn't a struct or pointer to a struct an
// error will be returned in the validate function.
func NewQueryStr(src interface{}, queryStr string) *Gator {
	g := New()
	m, err := url.ParseQuery(queryStr)
	if err != nil {
		err = fmt.Errorf("gator: couldn't parse QueryStr - %s", err)
		g.Add(errValidator{err: err})
		return g
	}

	objT, objV, err := getReflectInfo(src)
	if err != nil {
		g.Add(errValidator{err: err})
		return g
	}

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
	return g
}

// Add adds Validators to the Gator.
func (g *Gator) Add(v Validator, vals ...Validator) *Gator {
	g.vals = append(g.vals, v)
	g.vals = append(g.vals, vals...)
	return g
}

// Validate implements the Validator interface and returns an error if
// any of the Validators added return an error.
func (g *Gator) Validate() error {
	for _, v := range g.vals {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// A Field is a named value that is validated against a supplied Func.
type Field struct {
	name string
	src  interface{}
	f    Func
}

// NewField creates an initialized Field
func NewField(name string, src interface{}, f Func) *Field {
	return &Field{
		name: name,
		src:  src,
		f:    f,
	}
}

// Validate implements the Validator interface.  Field's Validate
// method calls the Func supplied during initialization.
func (f *Field) Validate() error {
	return f.f(f.name, f.src)
}

// RegisterStructTagToken registers custom tokens for gator struct tags.
func RegisterStructTagToken(token string, convFunc func(string) Func) {
	textToFuncMap[token] = convFunc
}

var (
	textToFuncMap = map[string]func(string) Func{}
)

func init() {
	RegisterStructTagToken("nonzero", func(s string) Func { return Nonzero() })
	RegisterStructTagToken("email", func(s string) Func { return Email() })
	RegisterStructTagToken("hexcolor", func(s string) Func { return HexColor() })
	RegisterStructTagToken("url", func(s string) Func { return URL() })
	RegisterStructTagToken("ip", func(s string) Func { return IP() })
	RegisterStructTagToken("alpha", func(s string) Func { return Alpha() })
	RegisterStructTagToken("num", func(s string) Func { return Num() })
	RegisterStructTagToken("alphanum", func(s string) Func { return AlphaNum() })
	RegisterStructTagToken("matches", func(s string) Func { return Matches(s) })
	RegisterStructTagToken("lat", func(s string) Func { return Lat() })
	RegisterStructTagToken("lon", func(s string) Func { return Lon() })
	RegisterStructTagToken("gt", func(s string) Func {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return Gt(n)
	})
	RegisterStructTagToken("gte", func(s string) Func {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return Gte(n)
	})
	RegisterStructTagToken("lt", func(s string) Func {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return Lt(n)
	})
	RegisterStructTagToken("lte", func(s string) Func {
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return Lte(n)
	})
	RegisterStructTagToken("in", func(s string) Func {
		list := strings.Split(s, ",")
		iList := []interface{}{}
		for _, i := range list {
			iList = append(iList, i)
		}
		return In(iList)
	})
	RegisterStructTagToken("notin", func(s string) Func {
		list := strings.Split(s, ",")
		iList := []interface{}{}
		for _, i := range list {
			iList = append(iList, i)
		}
		return NotIn(iList)
	})
	RegisterStructTagToken("len", func(s string) Func {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return Len(int(n))
	})
	RegisterStructTagToken("minlen", func(s string) Func {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return MinLen(int(n))
	})
	RegisterStructTagToken("maxlen", func(s string) Func {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return textErrorFunc(s, err)
		}
		return MaxLen(int(n))
	})
	RegisterStructTagToken("each", func(s string) Func {
		return Each(funcsFromTag(s)...)
	})
}

func textErrorFunc(s string, err error) Func {
	return func(name string, v interface{}) error {
		return fmt.Errorf("gator: tag for %s received parsing error - %s", name, err)
	}
}

func funcsFromTag(tag string) []Func {
	funcs := []Func{}
	sects := strings.Split(tag, "|")
	tSects := []string{}
	for _, s := range sects {
		tSects = append(tSects, strings.TrimSpace(s))
	}
	for _, s := range tSects {
		nCap := nonCaptureString(s)
		cap := captureString(s)
		for tag, f := range textToFuncMap {
			if nCap == tag {
				fn := f(cap)
				funcs = append(funcs, fn)
				break
			}
		}
	}
	return funcs
}

func nonCaptureString(s string) string {
	start := strings.Index(s, "(")
	if start == -1 {
		return s
	}
	return strings.TrimSpace(s[0:start])
}

func captureString(s string) string {
	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	if start == -1 || end == -1 || start > end {
		return ""
	}
	return strings.TrimSpace(s[start+1 : end])
}

type errValidator struct {
	err error
}

func (e errValidator) Validate() error {
	return e.err
}
