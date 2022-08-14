package go2ts

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const rn = `
`

type API struct {
	Method string
	Url    string
}

type Go2ts struct {
	types      []reflect.Type
	apis       map[reflect.Type]API
	values     map[reflect.Type]any
	makedTypes map[string]bool
}

func New() *Go2ts {
	return &Go2ts{
		values:     map[reflect.Type]any{},
		makedTypes: map[string]bool{},
		apis:       map[reflect.Type]API{},
	}
}

func (t *Go2ts) Add(v any) *Go2ts {
	ty := reflect.TypeOf(v)
	t.types = append(t.types, ty)
	t.values[ty] = v
	return t
}

func (t *Go2ts) AddApi(method, url string, v any) *Go2ts {
	ty := reflect.TypeOf(v)
	t.types = append(t.types, ty)
	t.apis[ty] = API{
		Method: method,
		Url:    url,
	}
	t.values[ty] = v
	return t
}

var all = ""

func (t *Go2ts) Format_all() string {
	if all != "" {
		return all
	}
	out := "// Auto create with go2ts" + rn
	out += "/* eslint-disable */" + rn
	for _, v := range t.types {
		out += rn + t.formatStruct(v)
	}
	all = out
	return out
}

func (t *Go2ts) Log() *Go2ts {
	fmt.Printf("%s", t.Format_all())
	return t
}

func (t *Go2ts) Write(filePath string) *Go2ts {
	codes := t.Format_all()
	// 内容一致就不更新
	if by, err := os.ReadFile(filePath); err == nil {
		if string(by) == codes {
			return t
		}
	}

	if err := os.WriteFile(filePath, []byte(codes), 0o777); err != nil {
		panic(err)
	}
	return t
}

func fixFuncType(subTypes []reflect.Type, input reflect.Type) ([]reflect.Type, string) {
	typeStr := input.String()
	relType := typeMap[typeStr]
	if relType == "" {
		if strings.Contains(typeStr, ".") && !strings.Contains(typeStr, "func") {
			subTypes = append(subTypes, input)
			relType = strings.Split(typeStr, ".")[1]
		} else {
			relType = "any"
		}
	}

	return subTypes, relType
}

func fixStructType(subTypes []reflect.Type, field reflect.StructField) ([]reflect.Type, string) {
	typeStr := field.Type.String()
	// typeStr := strings.ReplaceAll(strings.ReplaceAll(field.Type.String(), "[]", ""), ".", "")
	relType := typeMap[typeStr]
	tsType := field.Tag.Get("ts_type")
	if tsType != "" {
		relType = tsType
	}

	if relType == "" {
		if field.Type.Kind() == reflect.Struct {
			subTypes = append(subTypes, field.Type)
			relType = strings.Split(typeStr, ".")[1]
		} else if field.Type.Kind() == reflect.Slice {
			typ := field.Type.Elem()
			typStr := typ.String()
			if typStr == "[]interface {}" {
				relType = "any[]"
			} else {
				subTypes = append(subTypes, typ)
				if strings.Contains(typStr, ".") {
					relType = strings.Split(typStr, ".")[1] + "[]"
				} else {
					relType = typStr + "[]"
				}
			}

		} else {
			relType = "any"
		}
	}

	return subTypes, relType
}

func (t *Go2ts) formatApi(api API, ty reflect.Type) string {
	subTypes := []reflect.Type{}
	val := t.values[ty]
	point := reflect.ValueOf(val).Pointer()
	fnName := runtime.FuncForPC(point).Name()
	strs := strings.Split(fnName, "/")
	fnName = strs[len(strs)-1]
	fnName = strings.ToTitle(fnName[:1]) + fnName[1:]
	fnName = strings.Replace(fnName, ".", "", 1)
	out := fmt.Sprintf(`export const api%s = (`, fnName)
	numin := ty.NumIn()
	args := []string{}
	for i := 0; i < numin; i++ {
		var relType string
		subTypes, relType = fixFuncType(subTypes, ty.In(i))

		arg := strings.ToLower(relType[:1]) + relType[1:]
		if tsBaseType[arg] {
			arg = fmt.Sprintf("arg%d", i)
		}
		if i == numin-1 {
			out += fmt.Sprintf(`%s: %s`, arg, relType)
		} else {
			out += fmt.Sprintf(`%s: %s,`, arg, relType)
		}
		args = append(args, arg)

	}
	out += ")"

	if ty.NumOut() > 0 {
		var relType string
		subTypes, relType = fixFuncType(subTypes, ty.Out(0))
		out += fmt.Sprintf(`:Promise<%s> => {`, relType)
	} else {
		out += "any {"
	}
	out += fmt.Sprintf(`return (window as any).customFetch("%s", "%s", %s); }`, api.Method, api.Url, strings.Join(args, ", "))

	// 额外涉及到的类型
	for _, ty := range subTypes {
		out += rn + t.formatStruct(ty)
	}
	return out
}

func (t *Go2ts) formatFunc(ty reflect.Type) string {
	subTypes := []reflect.Type{}
	val := t.values[ty]
	point := reflect.ValueOf(val).Pointer()
	fnName := runtime.FuncForPC(point).Name()
	strs := strings.Split(fnName, "/")
	fnName = strs[len(strs)-1]
	fnName = strings.ToTitle(fnName[:1]) + fnName[1:]
	fnName = strings.Replace(fnName, ".", "", 1)
	out := fmt.Sprintf(`export type api%s = (`, fnName)
	numin := ty.NumIn()
	for i := 0; i < numin; i++ {
		var relType string
		subTypes, relType = fixFuncType(subTypes, ty.In(i))
		if relType == "" {
			continue
		}

		arg := strings.ToLower(relType[:1]) + relType[1:]
		if tsBaseType[arg] {
			arg = fmt.Sprintf("arg%d", i)
		}
		if i == numin-1 {
			out += fmt.Sprintf(`%s: %s`, arg, relType)
		} else {
			out += fmt.Sprintf(`%s: %s,`, arg, relType)
		}

	}
	out += ")=>"

	if ty.NumOut() > 0 {
		var relType string
		subTypes, relType = fixFuncType(subTypes, ty.Out(0))
		out += fmt.Sprintf(`Promise<%s>;`, relType)
	} else {
		out += "any;"
	}

	for _, ty := range subTypes {
		out += rn + t.formatStruct(ty)
	}
	return out
}

func (t *Go2ts) formatStruct(ty reflect.Type) string {
	tyStr := ty.String()
	if t.makedTypes[tyStr] {
		return ""
	}
	t.makedTypes[tyStr] = true
	out := fmt.Sprintf(`export interface %s {%s`, ty.Name(), rn)
	endFn := func() {}
	if ty.Kind() == reflect.Func {
		if api, ok := t.apis[ty]; ok {
			return t.formatApi(api, ty)
		}
		return t.formatFunc(ty)
	}
	if ty.Kind() == reflect.Slice {
		ty = ty.Elem()
		name := ty.Name()
		out = fmt.Sprintf(`interface arr_%s {%s`, name, rn)
		endFn = func() {
			out += fmt.Sprintf("%sexport type %s = arr_%s[]", rn, name, name)
		}
		// return fmt.Sprintf(`export type %s = any[]%s`, sliceName(ty), rn)
	}
	if ty.Kind() == reflect.Pointer {
		ty = ty.Elem()
	}
	subTypes := []reflect.Type{}

	var parse func(reflect.Type, map[string]bool) string
	parse = func(ty reflect.Type, appended map[string]bool) string {
		out := ""

		for i := 0; i < ty.NumField(); i++ {
			field := ty.Field(i)
			jsonStr := field.Tag.Get("json")
			json := strings.Split(jsonStr, ",")
			key := json[0]
			if appended[key] {
				continue
			}
			appended[key] = true
			opt := "?"
			validate := field.Tag.Get("validate")
			if strings.Contains(validate, "required") {
				opt = ""
			}
			if len(json) > 1 && strings.Contains(jsonStr, "omitempty") {
				opt = "?"
			}
			var relType string

			subTypes, relType = fixStructType(subTypes, field)
			if field.Anonymous {
				out += fmt.Sprintf(`  // Anonymous: %s%s: %s; %s`, key, opt, relType, rn)
				out += parse(field.Type, appended)
				out += fmt.Sprintf(`  // End: %s%s: %s; %s`, key, opt, relType, rn)
			} else if relType != "" {
				out += fmt.Sprintf(`  %s%s: %s; %s`, key, opt, relType, rn)
			}
		}

		return out
	}

	out += parse(ty, map[string]bool{})

	out += "}"
	for _, ty := range subTypes {
		out += rn + t.formatStruct(ty)
	}
	endFn()
	return out
}
