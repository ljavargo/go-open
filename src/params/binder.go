package params

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Binder struct {
	// Bind takes the name and type of the desired parameter and constructs it
	// from one or more values from Params.
	//
	// Example
	//
	// Request:
	//   url?id=123&ol[0]=1&ol[1]=2&ul[]=str&ul[]=array&user.Name=rob
	//
	// Action:
	//   Example.Action(id int, ol []int, ul []string, user User)
	//
	// Calls:
	//   Bind(params, "id", int): 123
	//   Bind(params, "ol", []int): {1, 2}
	//   Bind(params, "user", User): User{Name:"rob"}
	//
	// Note that only exported struct fields may be bound.
	Bind func(params *url.Values, name string, typ reflect.Type) reflect.Value
}

// An adapter for easily making one-key-value binders.
func ValueBinder(f func(value string, typ reflect.Type) reflect.Value) func(*url.Values, string, reflect.Type) reflect.Value {
	return func(params *url.Values, name string, typ reflect.Type) reflect.Value {
		p := *params
		vals, ok := p[name]
		if !ok || len(vals) == 0 {
			return reflect.Zero(typ)
		}

		return f(vals[0], typ)
	}
}

const (
	DEFAULT_DATE_FORMAT            = "2006-01-02"
	DEFAULT_DATETIME_FORMAT        = "2006-01-02 15:0"
	DEFAULT_DATETIME_FORMAT_SECOND = "2006-01-02 15:04:05"
)

var (
	// These are the lookups to find a Binder for any type of data.
	// The most specific binder found will be used (Type before Kind)
	TypeBinders = make(map[reflect.Type]Binder)
	KindBinders = make(map[reflect.Kind]Binder)

	// Applications can add custom time formats to this array, and they will be
	// automatically attempted when binding a time.Time.
	TimeFormats = []string{}

	IntBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			if len(val) == 0 {
				return reflect.Zero(typ)
			}
			intValue, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return reflect.Zero(typ)
			}
			pValue := reflect.New(typ)
			pValue.Elem().SetInt(intValue)
			return pValue.Elem()
		}),
	}

	UintBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			if len(val) == 0 {
				return reflect.Zero(typ)
			}
			uintValue, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return reflect.Zero(typ)
			}
			pValue := reflect.New(typ)
			pValue.Elem().SetUint(uintValue)
			return pValue.Elem()
		}),
	}

	FloatBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			if len(val) == 0 {
				return reflect.Zero(typ)
			}
			floatValue, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return reflect.Zero(typ)
			}
			pValue := reflect.New(typ)
			pValue.Elem().SetFloat(floatValue)
			return pValue.Elem()
		}),
	}

	StringBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(val)
		}),
	}

	// Booleans support a couple different value formats:
	// "true" and "false"
	// "on" and "" (a checkbox)
	// "1" and "0" (why not)
	BoolBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			v := strings.TrimSpace(strings.ToLower(val))
			switch v {
			case "true", "on", "1":
				return reflect.ValueOf(true)
			}
			// Return false by default.
			return reflect.ValueOf(false)
		}),
	}

	PointerBinder = Binder{
		Bind: func(params *url.Values, name string, typ reflect.Type) reflect.Value {
			//return nil if param is unset
			par := *params
			vals, ok := par[name]
			if !ok || len(vals) == 0 {
				return reflect.Zero(typ)
			}

			v := Bind(params, name, typ.Elem())

			p := reflect.New(v.Type()).Elem()
			p.Set(v)
			return p.Addr()
		},
	}

	TimeBinder = Binder{
		Bind: ValueBinder(func(val string, typ reflect.Type) reflect.Value {
			for _, f := range TimeFormats {
				if f == "" {
					continue
				}

				if strings.Contains(f, "07") || strings.Contains(f, "MST") {
					if r, err := time.Parse(f, val); err == nil {
						return reflect.ValueOf(r)
					}
				} else {
					if r, err := time.ParseInLocation(f, val, time.Local); err == nil {
						return reflect.ValueOf(r)
					}
				}
			}

			if unixInt, err := strconv.ParseInt(val, 10, 64); err == nil {
				return reflect.ValueOf(time.Unix(unixInt, 0))
			}

			return reflect.Zero(typ)
		}),
	}
)

// Sadly, the binder lookups can not be declared initialized -- that results in
// an "initialization loop" compile error.
func init() {
	KindBinders[reflect.Int] = IntBinder
	KindBinders[reflect.Int8] = IntBinder
	KindBinders[reflect.Int16] = IntBinder
	KindBinders[reflect.Int32] = IntBinder
	KindBinders[reflect.Int64] = IntBinder

	KindBinders[reflect.Uint] = UintBinder
	KindBinders[reflect.Uint8] = UintBinder
	KindBinders[reflect.Uint16] = UintBinder
	KindBinders[reflect.Uint32] = UintBinder
	KindBinders[reflect.Uint64] = UintBinder

	KindBinders[reflect.Float32] = FloatBinder
	KindBinders[reflect.Float64] = FloatBinder

	KindBinders[reflect.String] = StringBinder
	KindBinders[reflect.Bool] = BoolBinder
	KindBinders[reflect.Slice] = Binder{bindSlice}
	KindBinders[reflect.Struct] = Binder{bindStruct}
	KindBinders[reflect.Ptr] = PointerBinder

	TypeBinders[reflect.TypeOf(time.Time{})] = TimeBinder

	TimeFormats = append(TimeFormats, DEFAULT_DATE_FORMAT, DEFAULT_DATETIME_FORMAT, DEFAULT_DATETIME_FORMAT_SECOND, time.RFC3339)
}

// Used to keep track of the index for individual keyvalues.
type sliceValue struct {
	index int           // Index extracted from brackets.  If -1, no index was provided.
	value reflect.Value // the bound value for this slice element.
}

func checkSlice(key string, name string) bool {
	if strings.HasPrefix(key, name+".") {
		return true
	}

	return false
}

func bindSliceStruct(params *url.Values, name string, typ reflect.Type) reflect.Value {
	p := *params
	sliceValues := []sliceValue{}
	maxIndex := -1
	numNoIndex := 0
	for key, vals := range p {
		if len(strings.Split(key, ".")) < 3 || !strings.HasPrefix(key, name+".") {
			continue
		}

		keyWithoutPrefix := strings.TrimPrefix(key, name+".")

		ss := strings.Split(keyWithoutPrefix, ".")
		indexStr := ss[0]

		prefix := name + "." + indexStr
		index, ierr := strconv.Atoi(indexStr)
		if ierr != nil {
			fmt.Println("==== invalid slice st", ierr, indexStr)
			return reflect.Value{}
		}

		if index > maxIndex {
			maxIndex = index
		}

		numNoIndex += len(vals)
		sliceValues = append(sliceValues, sliceValue{
			index: index,
			value: Bind(params, prefix, typ.Elem()),
		})

	}

	resultArray := reflect.MakeSlice(typ, maxIndex+1, maxIndex+1+numNoIndex)
	for _, sv := range sliceValues {
		if sv.index != -1 {
			resultArray.Index(sv.index).Set(sv.value)
		} else {
			resultArray = reflect.Append(resultArray, sv.value)
		}
	}

	return resultArray

}

func bindSlice(params *url.Values, name string, typ reflect.Type) reflect.Value {
	// Collect an array of slice elements with their indexes (and the max index).
	maxIndex := -1
	numNoIndex := 0
	sliceValues := []sliceValue{}
	if typ.Elem().Kind() == reflect.Struct {
		return bindSliceStruct(params, name, typ)
	}

	// Factor out the common slice logic (between form values and files).
	processElement := func(key string, vals []string) {
		if !checkSlice(key, name) {
			return
		}

		ss := strings.Split(key, ".")
		index, err := strconv.Atoi(ss[1])
		if err != nil {
			return
		}

		if index > maxIndex {
			maxIndex = index
		}

		numNoIndex += len(vals)
		// Unindexed values can only be direct-bound.
		sliceValues = append(sliceValues, sliceValue{
			index: index,
			value: Bind(params, key, typ.Elem()),
		})
	}

	p := *params
	for key, vals := range p {
		processElement(key, vals)
	}

	resultArray := reflect.MakeSlice(typ, maxIndex+1, maxIndex+1+numNoIndex)
	for _, sv := range sliceValues {
		if sv.index != -1 {
			resultArray.Index(sv.index).Set(sv.value)
		} else {
			resultArray = reflect.Append(resultArray, sv.value)
		}
	}

	return resultArray
}

// Break on dots and brackets.
// e.g. bar => "bar", bar.baz => "bar", bar[0] => "bar"
func nextKey(key string) string {
	fieldLen := strings.IndexAny(key, ".[")
	if fieldLen == -1 {
		return key
	}
	return key[:fieldLen]
}

func getfieldByTag(v reflect.Value, tagName string) reflect.Value {
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typ2 := typ.Field(i)

		if typ2.PkgPath != "" && !typ2.Anonymous { // skip unexport
			continue
		}

		tag := typ2.Tag.Get("json")
		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = typ2.Name
		}

		if tag == tagName {
			return field
		}
	}

	return reflect.Value{}
}

func bindStruct(params *url.Values, name string, typ reflect.Type) reflect.Value {
	result := reflect.New(typ).Elem()
	fieldValues := make(map[string]reflect.Value)

	p := *params
	for key, _ := range p {
		if !strings.HasPrefix(key, name+".") {
			continue
		}

		// Get the name of the struct property.
		// Strip off the prefix. e.g. foo.bar.baz => bar.baz
		suffix := key[len(name)+1:]
		if ss := strings.Split(suffix, "."); len(ss) > 1 {
			suffix = ss[0]
		}

		fieldName := nextKey(suffix)
		fieldLen := len(fieldName)

		if _, ok := fieldValues[fieldName]; !ok {
			// Time to bind this field.  Get it and make sure we can set it.
			fieldValue := getfieldByTag(result, fieldName)
			if !fieldValue.IsValid() {
				continue
			}
			if !fieldValue.CanSet() {
				continue
			}

			boundVal := Bind(params, key[:len(name)+1+fieldLen], fieldValue.Type())
			if boundVal.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(boundVal.Convert(fieldValue.Type()))
			}

			fieldValues[fieldName] = boundVal
		}
	}

	return result
}

// Bind takes the name and type of the desired parameter and constructs it
// from one or more values from Params.
// Returns the zero value of the type upon any sort of failure.
func Bind(params *url.Values, name string, typ reflect.Type) reflect.Value {
	if binder, found := binderForType(typ); found {
		return binder.Bind(params, name, typ)
	}
	return reflect.Zero(typ)
}

func BindValue(val string, typ reflect.Type) reflect.Value {
	return Bind(&url.Values{"": {val}}, "", typ)
}

func binderForType(typ reflect.Type) (Binder, bool) {
	binder, ok := TypeBinders[typ]
	if !ok {
		binder, ok = KindBinders[typ.Kind()]
		if !ok {
			// WARN.Println("no binder for type:", typ)
			return Binder{}, false
		}
	}
	return binder, true
}
