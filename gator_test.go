package gator_test

import (
	"testing"

	"github.com/ShaleApps/gator"
)

type testStruct1 struct {
	Required string `gator:"nonzero"`
}

type testStruct2 struct {
	Email    string `gator:"email"`
	HexColor string `gator:"hexcolor"`
}

type testStruct3 struct {
	URL      string `gator:"url"`
	Username string `gator:"alphanum | minlen(5) | maxlen(10)"`
}

type testStruct4 struct {
	Int   int     `gator:"gt(18)"`
	Int2  int64   `gator:"gte(100)"`
	Float float64 `gator:"lt(19.9)"`
	Uint  uint    `gator:"lte(18)"`
}

type testStruct5 struct {
	IntSlice []int `gator:"each(gt(18))"`
}

type testStruct6 struct {
	Lat float64 `gator:"lat"`
	Lon float64 `gator:"lon"`
}

type testStruct7 struct {
	Eq1 string `gator:"eq(abc123)"`
	Eq2 string `gator:"eq(1)"`
	Eq3 int    `gator:"eq(1)"`
}

var validStructs = []interface{}{
	testStruct1{"a"},
	testStruct2{"loganjspears@gmail.com", "#ffffff"},
	testStruct2{"loganjspears@gmail.com", "#FFFFFF"},
	testStruct3{"http://www.google.com", "logan12345"},
	testStruct4{19, 400, -10.0000, 18},
	testStruct5{[]int{19, 20, 21}},
	testStruct6{0.0, 0.0},
	testStruct6{90.0, -180.0},
	testStruct7{"abc123", "1", 1},
}

var invalidStructs = []interface{}{
	testStruct1{Required: ""},
	testStruct2{"loganjspears@gmail", "#ffffff"},
	testStruct2{"loganjspears@gmail.com", "#fhffff"},
	testStruct3{"http://google", "logan12345"},
	testStruct3{"http://www.google.com", "log1"},
	testStruct3{"http://www.google.com", "logan100101001100101001"},
	testStruct4{14, 400, -10.0000, 18},
	testStruct4{19, 99, -10.0000, 18},
	testStruct4{19, 400, 19.90000001, 18},
	testStruct4{19, 400, -10.0000, 19},
	testStruct5{[]int{19, 20, 10}},
	testStruct6{90.1, -180.0},
	testStruct6{90.0, -180.1},
}

func TestStructs(t *testing.T) {
	for _, v := range validStructs {
		if err := gator.NewStruct(v).Validate(); err != nil {
			t.Errorf("struct should be valid: %+v %s", v, err)
		}
	}
	for _, v := range invalidStructs {
		if err := gator.NewStruct(v).Validate(); err == nil {
			t.Errorf("struct should be invalid: %+v", v)
		}
	}
}

var (
	validFields = []*gator.Field{
		gator.NewField("test", 1, gator.Nonzero()),
		gator.NewField("test", "abc", gator.Nonzero()),
		gator.NewField("test", testStruct1{"a"}, gator.Nonzero()),
		gator.NewField("test", &testStruct1{"a"}, gator.Nonzero()),
		gator.NewField("test", "test@example.com", gator.Email()),
		gator.NewField("test", "test+extension@reallyreallylongdomain.org", gator.Email()),
		gator.NewField("test", "lmy-us3r_n4m3", gator.Matches("^[a-z0-9_-]{3,16}$")),
		gator.NewField("test", "myp4ssw0rd", gator.Matches("^[a-z0-9_-]{6,18}$")),
		gator.NewField("test", 2, gator.Gt(1)),
		gator.NewField("test", 5, gator.Gt(4.9999999999999)),
		gator.NewField("test", -1, gator.Gt(-2)),
		gator.NewField("test", 1, gator.Lt(2)),
		gator.NewField("test", 4.9999999999999, gator.Lt(5)),
		gator.NewField("test", -2, gator.Lt(-1)),
		gator.NewField("test", 2, gator.Gte(2)),
		gator.NewField("test", 2, gator.Lte(2)),
		gator.NewField("test", 0, gator.Lat()),
		gator.NewField("test", 90.0, gator.Lat()),
		gator.NewField("test", -90.0, gator.Lat()),
		gator.NewField("test", 0, gator.Lon()),
		gator.NewField("test", 180.0, gator.Lon()),
		gator.NewField("test", -180.0, gator.Lon()),
		gator.NewField("test", "one", gator.In([]interface{}{"one", "two"})),
		gator.NewField("test", 1, gator.In([]interface{}{1, 2})),
		gator.NewField("test", "three", gator.NotIn([]interface{}{"one", "two"})),
		gator.NewField("test", 3, gator.NotIn([]interface{}{1, 2})),
		gator.NewField("test", "one", gator.NotIn([]interface{}{1, 2})),
		gator.NewField("test", "123456", gator.Len(6)),
		gator.NewField("test", []int{1, 2, 3}, gator.Len(3)),
		gator.NewField("test", "123456", gator.MinLen(5)),
		gator.NewField("test", "123456", gator.MinLen(6)),
		gator.NewField("test", []int{1, 2, 3}, gator.MinLen(2)),
		gator.NewField("test", []int{1, 2, 3}, gator.MinLen(3)),
		gator.NewField("test", "123456", gator.MaxLen(7)),
		gator.NewField("test", "123456", gator.MaxLen(6)),
		gator.NewField("test", []int{1, 2, 3}, gator.MaxLen(4)),
		gator.NewField("test", []int{1, 2, 3}, gator.MaxLen(3)),
		gator.NewField("test", []string{"1", "12", "123"}, gator.Each(gator.MaxLen(3))),
		gator.NewField("test", []int{1, 2, 3}, gator.Each(gator.Gt(0))),
		gator.NewField("test", 1, gator.Eq(1.0)),
		gator.NewField("test", "hello", gator.Eq("hello")),
		gator.NewField("test", "1", gator.Eq("1")),
	}

	invalidFields = []*gator.Field{
		gator.NewField("test", 0, gator.Nonzero()),
		gator.NewField("test", "", gator.Nonzero()),
		gator.NewField("test", testStruct1{}, gator.Nonzero()),
		gator.NewField("test", nil, gator.Nonzero()),
		gator.NewField("test", "test#example.com", gator.Email()),
		gator.NewField("test", "test @ reallyreallylongdomain.org", gator.Email()),
		gator.NewField("test", "@example.org", gator.Email()),
		gator.NewField("test", "th1s1s-wayt00_l0ngt0beausername", gator.Matches("^[a-z0-9_-]{3,16}$")),
		gator.NewField("test", "mypa$$w0rd", gator.Matches("^[a-z0-9_-]{6,18}$")),
		gator.NewField("test", 1, gator.Matches("^[a-z0-9_-]{6,18}$")),
		gator.NewField("test", 1, gator.Gt(2)),
		gator.NewField("test", 4.9999999999999, gator.Gt(5)),
		gator.NewField("test", -2, gator.Gt(-1)),
		gator.NewField("test", 2, gator.Gt(2)),
		gator.NewField("test", 2, gator.Lt(1)),
		gator.NewField("test", 5, gator.Lt(4.9999999999999)),
		gator.NewField("test", -1, gator.Lt(-2)),
		gator.NewField("test", 0, gator.Lt(0)),
		gator.NewField("test", 90.1, gator.Lat()),
		gator.NewField("test", -90.1, gator.Lat()),
		gator.NewField("test", 180.1, gator.Lon()),
		gator.NewField("test", -180.1, gator.Lon()),
		gator.NewField("test", "three", gator.In([]interface{}{"one", "two"})),
		gator.NewField("test", 3, gator.In([]interface{}{1, 2})),
		gator.NewField("test", "one", gator.In([]interface{}{1, 2})),
		gator.NewField("test", "one", gator.NotIn([]interface{}{"one", "two"})),
		gator.NewField("test", 1, gator.NotIn([]interface{}{1, 2})),
		gator.NewField("test", "123456", gator.Len(7)),
		gator.NewField("test", []int{1, 2, 3}, gator.Len(2)),
		gator.NewField("test", "123456", gator.MinLen(7)),
		gator.NewField("test", []int{1, 2, 3}, gator.MinLen(4)),
		gator.NewField("test", "123456", gator.MaxLen(5)),
		gator.NewField("test", []int{1, 2, 3}, gator.MaxLen(-1)),
		gator.NewField("test", []string{"1", "12", "123"}, gator.Each(gator.MaxLen(1))),
		gator.NewField("test", []int{1, 2, 3}, gator.Each(gator.Lt(3))),
		gator.NewField("test", -1, gator.Eq(1.0)),
		gator.NewField("test", "hello", gator.Eq("hell0")),
	}
)

func TestGator(t *testing.T) {
	for _, f := range validFields {
		if err := gator.New().Add(f).Validate(); err != nil {
			t.Errorf("%+v should have been valid, but produced an error: %s", f, err)
		}
	}
	for _, f := range invalidFields {
		if err := gator.New().Add(f).Validate(); err == nil {
			t.Errorf("%+v should have been invalid, but failed to produce an error", f)
		}
	}
}

type queryTest struct {
	QueryStr string
	Src      interface{}
}

var (
	validQueries = []*queryTest{
		&queryTest{
			QueryStr: `URL=url&Username=alphanum|minlen(5)|maxlen(10)`,
			Src: &testStruct3{
				URL:      "https://news.ycombinator.com",
				Username: "hello1",
			},
		},
	}

	invalidQueries = []*queryTest{
		&queryTest{
			QueryStr: `URL=url&Username=alphanum|minlen(5)|maxlen(10)`,
			Src: &testStruct3{
				URL:      "https//news.ycombinator.com",
				Username: "hello1",
			},
		},
		&queryTest{
			QueryStr: `URL=url&Username=alphanum|minlen(5)|maxlen(10)`,
			Src: &testStruct3{
				URL:      "https://news.ycombinator.com",
				Username: "hello",
			},
		},
		&queryTest{
			QueryStr: `URL=minlen(sss)`,
			Src: &testStruct3{
				URL:      "https//news.ycombinator.com",
				Username: "hello1",
			},
		},
	}
)

func TestQueryStr(t *testing.T) {
	for _, queryTest := range validQueries {
		if err := gator.NewQueryStr(queryTest.Src, queryTest.QueryStr).Validate(); err != nil {
			t.Errorf("%+v should have been valid, but produced an error: %s", queryTest.Src, err)
		}
	}

	for _, queryTest := range invalidQueries {
		if err := gator.NewQueryStr(queryTest.Src, queryTest.QueryStr).Validate(); err == nil {
			t.Errorf("%+v should have been invalid, but failed to produce an error", queryTest.Src)
		}
	}
}

func TestRegister(t *testing.T) {
	gator.RegisterStructTagToken("pword", func(s string) gator.Func {
		return gator.Matches(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{4,8}$`)
	})
	type User struct {
		Email    string `gator:”email”`
		Password string `gator:"pword"`
	}
	u := &User{
		Email:    "gator@example.com",
		Password: "ASDF12345",
	}
	if err := gator.NewStruct(u).Validate(); err == nil {
		t.Errorf("%+v should have been invalid, but failed to produce an error", u)
	}
}
