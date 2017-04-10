package params

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	assert "github.com/stretchr/testify.v2/require"
)

func TestPositiveBindValuesToStruct(t *testing.T) {
	type aliasString string

	type embeddedStruct struct {
		Embedded string `json:"embedded"`
	}

	type St2 struct {
		Hello string `json:"hello"`
	}

	type Testst struct {
		Hellow string      `json:"hellow"`
		Hello  aliasString `json:"hello"`
		He     []St2       `json:"he"`
		Nice   St2         `json:"nice"`
	}

	type input struct {
		// embeded field
		embeddedStruct

		// builtin type
		Int64Min  int64   `json:"int64min"`
		Int64Max  int64   `json:"int64max"`
		Uint64Min uint64  `json:"uint64min"`
		Uint64Max uint64  `json:"uint64max"`
		Float64   float64 `json:"float64"`
		String    string  `json:"string"`
		BoolTrue  bool    `json:"booltrue"`
		BoolOn    bool    `json:"boolon"`
		Bool1     bool    `json:"bool1"`
		BoolFalse bool    `json:"boolfalse"`

		// time
		Date        time.Time `json:"date"`
		DateTime    time.Time `json:"datetime"`
		TimeUnix    time.Time `json:"timeunix"`
		TimeRFC3339 time.Time `json:"timerfc3339"`

		// alias type
		AliasString aliasString `json:"aliasstring"`
		Hello       []string    `json:"hello"`
		Testst      []Testst    `json:"testst"`
		Test2t      Testst      `json:"tests"`
	}

	paramsIn := url.Values{
		"embedded": []string{"helloworld"},

		"int64min":  []string{"-9223372036854775808"},
		"int64max":  []string{"9223372036854775807"},
		"uint64min": []string{"0"},
		"uint64max": []string{"18446744073709551615"},
		"float64":   []string{"3.1415926"},
		"string":    []string{"abcdefg"},
		"booltrue":  []string{"true"},
		"boolon":    []string{"on"},
		"bool1":     []string{"1"},
		"boolfalse": []string{"false"},
		"fuck":      []string{"fuck"},

		"date":        []string{"2006-01-02"},
		"datetime":    []string{"2006-01-02 15:04:05"},
		"timeunix":    []string{"1136214245"},
		"timerfc3339": []string{"2006-01-02T15:04:05Z08:00"},

		"aliasstring":         []string{"alphago"},
		"hello.0":             []string{"hello0"},
		"hello.1":             []string{"hello1"},
		"hello.2":             []string{"hello2"},
		"tests.hellow":        []string{"hellow"},
		"tests.hello":         []string{"hello"},
		"testst.0.hellow":     []string{"helloworld"},
		"testst.0.hello":      []string{"helloworld"},
		"testst.1.hellow":     []string{"完"},
		"testst.2.hello":      []string{"helloworld"},
		"testst.2.he.0.hello": []string{"helloworld"},
		"testst.2.nice.hello": []string{"helloworld"},
	}

	var structActual input

	date, _ := time.ParseInLocation("2006-01-02", "2006-01-02", time.Local)
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2006-01-02 15:04:05", time.Local)
	timeunix := time.Unix(1136214245, 0)
	timerfc3339, _ := time.ParseInLocation(time.RFC3339, "2006-01-02T15:04:05Z08:00", time.Local)

	structExpected := input{
		embeddedStruct: embeddedStruct{Embedded: "helloworld"},

		Int64Min:  -9223372036854775808,
		Int64Max:  9223372036854775807,
		Uint64Min: 0,
		Uint64Max: 18446744073709551615,
		Float64:   3.1415926,
		String:    "abcdefg",
		BoolTrue:  true,
		BoolOn:    true,
		Bool1:     true,
		BoolFalse: false,

		Date:        date,
		DateTime:    datetime,
		TimeUnix:    timeunix,
		TimeRFC3339: timerfc3339,

		AliasString: aliasString("alphago"),
		Hello:       []string{"hello0", "hello1", "hello2"},
		Test2t: Testst{
			Hellow: "hellow",
			Hello:  "hello",
		},
		Testst: []Testst{
			Testst{
				Hello:  "helloworld",
				Hellow: "helloworld",
			},
			Testst{
				Hellow: "完",
			},
			Testst{
				Hello: "helloworld",
				He:    []St2{St2{Hello: "helloworld"}},
				Nice:  St2{Hello: "helloworld"},
			},
		},
	}

	req, _ := http.NewRequest("GET", "http://localhost?"+paramsIn.Encode(), nil)

	BindValuesToStruct(&structActual, req)

	assert.Equal(t, structExpected.Embedded, structActual.Embedded, "embeded testing")
	assert.Equal(t, structExpected.Int64Min, structActual.Int64Min, "int64 min testing")
	assert.Equal(t, structExpected.Int64Max, structActual.Int64Max, "int64 max testing")
	assert.Equal(t, structExpected.Uint64Min, structActual.Uint64Min, "uint64 min testing")
	assert.Equal(t, structExpected.Uint64Max, structActual.Uint64Max, "uint64 max testing")
	assert.Equal(t, structExpected.Float64, structActual.Float64, "float64 testing")
	assert.Equal(t, structExpected.String, structActual.String, "string testing")
	assert.Equal(t, structExpected.BoolTrue, structActual.BoolTrue, "bool 'true' testing")
	assert.Equal(t, structExpected.BoolOn, structActual.BoolOn, "bool 'on' testing")
	assert.Equal(t, structExpected.Bool1, structActual.Bool1, "bool '1' testing")
	assert.Equal(t, structExpected.BoolFalse, structActual.BoolFalse, "bool 'false' testing")
	assert.Equal(t, structExpected.Date.Unix(), structActual.Date.Unix(), "date testing")
	assert.Equal(t, structExpected.DateTime.Unix(), structActual.DateTime.Unix(), "datetime testing")
	assert.Equal(t, structExpected.TimeUnix.Unix(), structActual.TimeUnix.Unix(), "unix time testing")
	assert.Equal(t, structExpected.TimeRFC3339.Unix(), structActual.TimeRFC3339.Unix(), "RFC3339 time testing")
	assert.Equal(t, structExpected.AliasString, structActual.AliasString, "alias testing")
	assert.Equal(t, structExpected.Hello[0], structActual.Hello[0], "slice")
	assert.Equal(t, structExpected.Hello[1], structActual.Hello[1], "slice")
	assert.Equal(t, structExpected.Hello[2], structActual.Hello[2], "slice")
	assert.Equal(t, structExpected.Test2t.Hello, structActual.Test2t.Hello, "st")
	assert.Equal(t, structExpected.Test2t.Hellow, structActual.Test2t.Hellow, "st")
	assert.Equal(t, structExpected.Testst[0].Hello, structActual.Testst[0].Hello, "st slice")
	assert.Equal(t, structExpected.Testst[0].Hellow, structActual.Testst[0].Hellow, "st slice")
	assert.Equal(t, structExpected.Testst[1].Hellow, structActual.Testst[1].Hellow, "st slice")
	assert.Equal(t, structExpected.Testst[2].Hello, structActual.Testst[2].Hello, "st slice")
	assert.Equal(t, structExpected.Testst[2].Nice.Hello, structActual.Testst[2].Nice.Hello, "st slice")
	assert.Equal(t, structExpected.Testst[2].He[0].Hello, structActual.Testst[2].He[0].Hello, "st slice")

}
