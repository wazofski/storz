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

	// TODO add imports

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

	methods := []string{}
	fields := []string{}
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

func writeStruct(b *strings.Builder, name string, fields []string) {
	const wrapperPrefix string = "Wrapper"

	write(b, fmt.Sprintf("type %s%s struct {", name, wrapperPrefix), 0)
	for _, f := range fields {
		write(b, f, 1)
	}
	write(b, "}", 0)
	endl(b)

	for _, f := range fields {
		tok := strings.Split(f, " ")
		nm := tok[0]
		tp := tok[1]
		pp := capitalize(nm)

		writeFunction(b,
			fmt.Sprintf("%s() %s", pp, tp),
			fmt.Sprintf("(o %s%s)", name, wrapperPrefix),
			[]string{fmt.Sprintf("return o.%s", nm)})

		writeFunction(b,
			fmt.Sprintf("Set%s(v %s)", pp, tp),
			fmt.Sprintf("(o %s%s)", name, wrapperPrefix),
			[]string{fmt.Sprintf("o.%s = v", nm)})
	}
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
