package mgen

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"strings"
)

func Generate(model string) error {
	structs, resources := loadModel(model)

	imports := []string{
		// "errors",
		// "log",
		// "strings",
		"fmt",
		"encoding/json",
		"github.com/wazofski/store/utils",
		"github.com/wazofski/store",
	}

	var b strings.Builder
	b.WriteString(render("mgen/templates/imports.gotext", imports))
	b.WriteString(compileResources(resources))
	b.WriteString(compileStructs(structs))

	str := strings.ReplaceAll(b.String(), "&#34;", "\"")
	res, err := format.Source([]byte(str))

	if err != nil {
		log.Println(err)
		res = []byte(str)
	}

	return exportFile("generated", "objects.go", string(res))
}

type _Interface struct {
	Name       string
	Methods    []string
	Implements []string
}

func compileResources(resources []_Resource) string {
	var b strings.Builder

	for _, r := range resources {
		props := []_Prop{
			{
				Prop:    "Meta",
				Type:    "store.Meta",
				Json:    "metadata",
				Default: fmt.Sprintf("store.MetaFactory(\"%s\")", r.Name),
			},
		}

		if len(r.Spec) > 0 {
			props = append(props,
				_Prop{
					Prop: "Spec",
					Type: r.Spec,
					Json: "spec",
				})
		}

		if len(r.Status) > 0 {
			props = append(props,
				_Prop{
					Prop: "Status",
					Type: r.Status,
					Json: "status",
				})
		}

		s := _Struct{
			Name:   r.Name,
			Props:  props,
			Embeds: []string{},
			Implements: []string{
				"store.Object",
			},
		}

		b.WriteString(compileStruct(s))
		b.WriteString(render("mgen/templates/meta.gotext", r))
		b.WriteString(render("mgen/templates/clone.gotext", s))
	}

	b.WriteString(render("mgen/templates/schema.gotext", resources))

	return b.String()
}

func compileStructs(structs []_Struct) string {
	var b strings.Builder

	for _, s := range structs {
		b.WriteString(compileStruct(s))
	}

	return b.String()
}

type _Tuple struct {
	A string
	B string
}

func compileStruct(s _Struct) string {
	var b strings.Builder
	methods := []string{}

	s.Props = addDefaultPropValues(s.Props)

	for _, p := range s.Props {
		if p.Prop != "Meta" {
			methods = append(methods,
				fmt.Sprintf("%s() %s", p.Prop, p.Type))
		}

		if p.Prop != "Meta" && p.Prop != "Spec" && p.Prop != "Status" {
			methods = append(methods,
				fmt.Sprintf("Set%s(v %s)", p.Prop, p.Type))
		}

		if p.Prop == "Spec" {
			b.WriteString(
				render("mgen/templates/specinternal.gotext",
					_Tuple{A: s.Name, B: p.Type}))
		}
	}

	impl := append(s.Implements, "json.Unmarshaler")

	b.WriteString(render("mgen/templates/interface.gotext", _Interface{
		Name:       s.Name,
		Methods:    methods,
		Implements: impl,
	}))

	b.WriteString(render("mgen/templates/structure.gotext", s))
	b.WriteString(render("mgen/templates/unmarshall.gotext", s))

	return b.String()
}

func render(path string, data interface{}) string {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)

	if err != nil {
		log.Fatalln(err)
	}

	return buf.String()
}

func addDefaultPropValues(props []_Prop) []_Prop {
	res := []_Prop{}

	for _, p := range props {
		if len(p.Default) > 0 {
			res = append(res, p)
			continue
		}

		res = append(res, _Prop{
			Prop:    p.Prop,
			Json:    p.Json,
			Type:    p.Type,
			Default: typeDefault(p.Type),
		})
	}

	return res
}

func typeDefault(tp string) string {
	if strings.HasPrefix(tp, "[]") {
		return fmt.Sprintf("%s {}", tp)
	}
	if strings.HasPrefix(tp, "map") {
		return fmt.Sprintf("make(%s)", tp)
	}

	switch tp {
	case "string":
		return "fmt.Sprint()"
	case "bool":
		return "false"
	case "int":
		return "0"
	case "float":
		return "0"
	default:
		return fmt.Sprintf("%sFactory()", tp)
	}
}

func (u _Prop) IsMap() bool {
	if len(u.Type) < 3 {
		return false
	}

	return u.Type[:3] == "map"
}

func (u _Prop) IsArray() bool {
	if len(u.Type) < 2 {
		return false
	}

	return u.Type[:2] == "[]"
}

func (u _Prop) StrippedType() string {
	if u.IsMap() {
		return u.Type[strings.LastIndex(u.Type, "]")+1:]
	}
	if u.IsArray() {
		return u.Type[2:]
	}
	return u.Type
}

func (u _Prop) StrippedDefault() string {
	return typeDefault(u.StrippedType())
}
