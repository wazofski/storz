package mgen

import (
	"fmt"
	"strings"
)

func Generate(path string, dest string) error {
	structs, resources := loadModel(path)

	var b strings.Builder

	packge := "generated"
	write(&b, fmt.Sprintf("package %s", packge), 0)
	endl(&b)

	imports := []string{
		"log",
		"encoding/json",
		"github.com/wazofski/store",
	}

	write(&b, compileImports(imports), 0)
	write(&b, compileStructs(structs), 0)
	write(&b, compileResources(resources), 0)

	fullDest := fmt.Sprintf("%s/%s/", dest, packge)
	return exportFile(fullDest, "objects.go", b.String())
}

func compileStructs(structs map[string]_Struct) string {
	var b strings.Builder

	for _, s := range structs {
		write(&b, compileStruct(s), 0)
	}

	return b.String()
}

func compileImports(imports []string) string {
	var b strings.Builder

	write(&b, "import (", 0)
	for _, s := range imports {
		write(&b, fmt.Sprintf("\"%s\"", s), 1)
	}
	write(&b, ")", 0)

	return b.String()
}

func compileResources(resources map[string]_Resource) string {
	var b strings.Builder

	for _, s := range resources {
		write(&b, compileResource(s), 0)
	}

	return b.String()
}

func compileStruct(str _Struct) string {
	var b strings.Builder

	methods := []string{}
	fields := []string{}
	for _, p := range str.Props {
		pp := capitalize(p.Prop)
		methods = append(methods, fmt.Sprintf("%s() %s", pp, p.Type))
		methods = append(methods, fmt.Sprintf("Set%s(v %s)", pp, p.Type))
		fields = append(fields, fmt.Sprintf("%s %s", p.Prop, p.Type))
	}

	writeInterface(&b, str.Name, methods)
	writeStruct(&b, str.Name, fields)

	return b.String()
}

func compileResource(str _Resource) string {
	var b strings.Builder

	methods := []string{"store.Object"}
	fields := []string{"store.ObjectWrapper"}
	if len(str.Spec) > 0 {
		methods = append(methods, fmt.Sprintf("Spec() %s", str.Spec))
		fields = append(fields, fmt.Sprintf("spec %s", str.Spec))
	}
	if len(str.Status) > 0 {
		methods = append(methods, fmt.Sprintf("Status() %s", str.Status))
		fields = append(fields, fmt.Sprintf("status %s", str.Status))
	}
	writeInterface(&b, str.Name, methods)
	writeStruct(&b, str.Name, fields)

	lines := []string{"return o.ObjectWrapper.Metadata"}
	writeFunction(&b,
		"Metadata() store.Meta",
		fmt.Sprintf("(o *%s%s)", str.Name, wrapperSuffix),
		lines)

	lines = []string{
		"res, err := json.Marshal(*o)",
		"if err != nil {",
		"	log.Fatalln(err)",
		"}",
		"return res",
	}
	writeFunction(&b,
		"Serialize() []byte",
		fmt.Sprintf("(o *%s%s)", str.Name, wrapperSuffix),
		lines)

	return b.String()
}

func writeInterface(b *strings.Builder, name string, methods []string) {
	write(b, fmt.Sprintf("type %s interface {", name), 0)
	for _, m := range methods {
		write(b, m, 1)
	}
	write(b, "}", 0)
	endl(b)
}

const wrapperSuffix string = "Wrapper"
const factorySuffix string = "Factory"

func writeStruct(b *strings.Builder, name string, fields []string) {
	write(b, fmt.Sprintf("type %s%s struct {", name, wrapperSuffix), 0)

	for _, f := range fields {
		write(b, f, 1)
	}

	write(b, "}", 0)
	endl(b)

	lines := []string{
		fmt.Sprintf("return &%s%s {", name, wrapperSuffix),
	}

	for _, f := range fields {
		tok := strings.Split(f, " ")
		nm := tok[0]
		if strings.HasSuffix(nm, "ObjectWrapper") {
			lines = append(lines, fmt.Sprintf("ObjectWrapper: store.ObjectWrapperFactory(\"%s\"),", name))
			continue
		}
		tp := tok[1]
		pp := capitalize(nm)
		lines = append(lines, fmt.Sprintf("%s: %s,", nm, typeDefault(tp)))

		writeFunction(b,
			fmt.Sprintf("%s() %s", pp, tp),
			fmt.Sprintf("(o *%s%s)", name, wrapperSuffix),
			[]string{fmt.Sprintf("return o.%s", nm)})

		if nm != "spec" && nm != "status" {
			writeFunction(b,
				fmt.Sprintf("Set%s(v %s)", pp, tp),
				fmt.Sprintf("(o *%s%s)", name, wrapperSuffix),
				[]string{fmt.Sprintf("o.%s = v", nm)})
		}
	}

	endl(b)

	lines = append(lines, "}")
	writeFunction(b,
		fmt.Sprintf("%s%s() %s", name, factorySuffix, name),
		"",
		lines)
}

func writeFunction(
	b *strings.Builder,
	name string,
	instance string,
	lines []string) {

	write(b,
		fmt.Sprintf(
			"func %s %s {", instance, name),
		0)

	for _, l := range lines {
		write(b, l, 1)
	}

	write(b, "}", 0)
	endl(b)
}

func typeDefault(tp string) string {
	if strings.HasPrefix(tp, "[]") {
		return fmt.Sprintf("%s {}", tp)
	}
	if strings.HasPrefix(tp, "map") {
		return fmt.Sprintf("make(%s)", tp)
	}

	if tp == "string" {
		return "\"\""
	}
	if tp == "bool" {
		return "false"
	}
	if tp == "int" {
		return "0"
	}
	if tp == "float" {
		return "0"
	}

	return fmt.Sprintf("%s%s()", tp, factorySuffix)
}
