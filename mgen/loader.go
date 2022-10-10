package mgen

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

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
	Status     string       `yaml:"status,omitempty"`
	ApiMethods []_ApiMethod `yaml:"apimethods,omitempty"`
	Props      []_Prop      `yaml:"props,omitempty"`
}

type _Model struct {
	Types []_Type `yaml:"types"`
}

type _Struct struct {
	Name  string  `yaml:"name"`
	Props []_Prop `yaml:"props,omitempty"`
}

type _Resource struct {
	Name       string       `yaml:"name"`
	Spec       string       `yaml:"spec,omitempty"`
	Status     string       `yaml:"status,omitempty"`
	ApiMethods []_ApiMethod `yaml:"apimethods,omitempty"`
}

func loadModel(path string) (map[string]_Struct, map[string]_Resource) {
	yamls := findYamlFiles(path)
	structs := make(map[string]_Struct)
	resources := make(map[string]_Resource)

	for _, y := range yamls {
		model, err := readModel(y)
		if err != nil {
			log.Fatal(err)
		}

		for _, m := range model.Types {
			if m.Kind == "Struct" {
				structs[m.Name] = _Struct{
					Name:  m.Name,
					Props: m.Props,
				}
				continue
			}
			if m.Kind == "Resource" {
				resources[m.Name] = _Resource{
					Name:       m.Name,
					Spec:       m.Spec,
					Status:     m.Status,
					ApiMethods: m.ApiMethods,
				}
				continue
			}
		}
	}

	return structs, resources
}

func readModel(path string) (*_Model, error) {
	fmt.Printf("reading model %s ", path)
	fmt.Println()

	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := _Model{}
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		return nil, err
	}

	// for _, t := range data.Types {
	// 	fmt.Printf("%s", t)
	// 	fmt.Println()
	// }

	return &data, nil
}
