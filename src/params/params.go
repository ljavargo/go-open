package params

import (
	"net/http"
	"net/url"
	"reflect"
)

func BindValuesToStruct(dest interface{}, req *http.Request, returnvalue ...bool) reflect.Value {
	pointerMap := make(map[uintptr]bool)
	val := reflect.ValueOf(dest)
	elm := reflect.Indirect(val)
	if val.Kind() != reflect.Ptr && elm.Kind() != reflect.Struct {
		panic("need ptr of struct")
	}

	p := req.URL.Query()

	if len(returnvalue) > 0 {
		return bindValuesToStructWhitValue(elm, pointerMap, &p)
	} else {
		bindValuesToStruct(elm, pointerMap, &p)
	}

	return reflect.Value{}
}

func bindValuesToStruct(elm reflect.Value, pointerMap map[uintptr]bool, p *url.Values) {
	typ := elm.Type()
	bind(elm, typ, pointerMap, p)
}

func bindValuesToStructWhitValue(elm reflect.Value, pointerMap map[uintptr]bool, p *url.Values) reflect.Value {
	typ := elm.Type()
	result := reflect.New(typ).Elem()
	bind(result, typ, pointerMap, p)
	return result
}

func bind(elm reflect.Value, typ reflect.Type, pointerMap map[uintptr]bool, p *url.Values) {
	for i := 0; i < elm.NumField(); i++ {
		field := elm.Field(i)
		ftyp := typ.Field(i)

		if ftyp.PkgPath != "" && !ftyp.Anonymous { // skip unexport
			continue
		}

		tag := ftyp.Tag.Get("json")
		if tag == "-" {
			continue
		}

		name := tag
		if name == "" {
			name = ftyp.Name
		}

		// struct recursion
		if ftyp.Anonymous {
			var inf reflect.Value

			if field.Kind() == reflect.Ptr {
				pointer := field.Pointer()
				if pointerMap[pointer] {
					continue
				}

				subElm := reflect.Indirect(field)
				if subElm.Kind() != reflect.Struct {
					continue
				}
				inf = subElm

				// save pointer
				pointerMap[pointer] = true
			} else if field.Kind() == reflect.Struct {
				inf = field
			} else {
				continue
			}

			bindValuesToStruct(inf, pointerMap, p)
		} else {
			paramValue := Bind(p, name, field.Type())
			if paramValue.Type().ConvertibleTo(field.Type()) {
				field.Set(paramValue.Convert(field.Type()))
			}
		}
	}
}
