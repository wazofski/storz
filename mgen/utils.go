package mgen

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func yamlFiles(path string) []string {
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

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func decapitalize(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func exportFile(targetDir string, name string, content string) error {
	targetDir = strings.ReplaceAll(targetDir, "//", "/")
	os.RemoveAll(targetDir)
	os.Mkdir(targetDir, 0755)

	targetFile := targetDir + name

	log.Printf("exporting file %s", targetFile)

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(content)

	return nil
}
