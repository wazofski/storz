package mgen

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func findYamlFiles(path string) []string {
	res := []string{}

	libRegEx, e := regexp.Compile(`^\w+\.(?:yaml|yml)$`)
	if e != nil {
		log.Fatal(e)
	}

	e = filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err == nil && libRegEx.MatchString(info.Name()) {
				res = append(res, path)
			}
			return nil
		})

	if e != nil {
		log.Fatal(e)
	}

	return res
}

func write(b *strings.Builder, s string, indentation int) {
	const tabs string = "    "

	for i := 0; i < indentation; i++ {
		b.WriteString(tabs)
	}

	b.WriteString(s)
	endl(b)
}

func endl(b *strings.Builder) {
	b.WriteString(`
`)
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
