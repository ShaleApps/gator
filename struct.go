package gator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	structTagKey = "gator"
)

type Struct struct {
	src interface{}
}

func NewStruct(src interface{}) *Struct {
	return &Struct{src: src}
}

func (s *Struct) Validate() error {
	if err := isStructOrStructPtr(s.src); err != nil {
		return err
	}
	objT := reflect.TypeOf(s.src)
	objV := reflect.ValueOf(s.src)

	g := New()
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
	return g.Validate()
}

var (
	textToFuncMap = map[string]func(string) Func{
		"nonzero": func(s string) Func {
			return Nonzero
		},
		"email": func(s string) Func {
			return Email
		},
		"hexcolor": func(s string) Func {
			return HexColor
		},
		"url": func(s string) Func {
			return URL
		},
		"ip": func(s string) Func {
			return IP
		},
		"alpha": func(s string) Func {
			return Alpha
		},
		"num": func(s string) Func {
			return Num
		},
		"alphanum": func(s string) Func {
			return AlphaNum
		},
		"gt": func(s string) Func {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return Gt(n)
		},
		"gte": func(s string) Func {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return Gte(n)
		},
		"lt": func(s string) Func {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return Lt(n)
		},
		"lte": func(s string) Func {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return Lte(n)
		},
		"lat": func(s string) Func {
			return Lat
		},
		"lon": func(s string) Func {
			return Lon
		},
		"in": func(s string) Func {
			list := strings.Split(s, ",")
			iList := []interface{}{}
			for _, i := range list {
				iList = append(iList, i)
			}
			return In(iList)
		},
		"notin": func(s string) Func {
			list := strings.Split(s, ",")
			iList := []interface{}{}
			for _, i := range list {
				iList = append(iList, i)
			}
			return NotIn(iList)
		},
		"matches": func(s string) Func {
			return Matches(s)
		},
		"len": func(s string) Func {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return Len(int(n))
		},
		"minlen": func(s string) Func {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return MinLen(int(n))
		},
		"maxlen": func(s string) Func {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return textErrorFunc(s, err)
			}
			return MaxLen(int(n))
		},
	}
)

func init() {
	// can't put in textToFuncMap initialization because of self reference
	textToFuncMap["each"] = func(s string) Func {
		return Each(funcsFromTag(s)...)
	}
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
