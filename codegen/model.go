package codegen

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func ProcessModel(path string) error {
	yamls := findYamlFiles(path)
	for _, y := range yamls {
		err := readModel(y)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

type _ApiMethod struct {
	ApiMethod string `yaml:"apimethod,omitempty"`
}

type _Prop struct {
	Prop string `yaml:"prop"`
	Type string `yaml:"type"`
}

type _Type struct {
	Name       string       `yaml:"name"`
	Kind       string       `yaml:"kind,omitempty"`
	Spec       string       `yaml:"spec,omitempty"`
	ApiMethods []_ApiMethod `yaml:"apimethods,omitempty"`
	Props      []_Prop      `yaml:"props,omitempty"`
}

type _Model struct {
	Types []_Type `yaml:"types"`
}

func readModel(path string) error {
	fmt.Printf("reading model %s ", path)
	fmt.Println()

	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	data := _Model{}
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		return err
	}

	for i, d := range data.Types {
		fmt.Printf("row %d %s", i, d)
		fmt.Println()
	}

	return nil
}
