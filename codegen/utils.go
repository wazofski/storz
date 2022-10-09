package codegen

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
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
