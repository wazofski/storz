package mgen

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type _ApiMethod struct {
	ApiMethod string `yaml:"apimethod,omitempty"`
}

type _Prop struct {
	Prop    string `yaml:"prop"`
	Type    string `yaml:"type"`
	Json    string
	Default string
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
	Name       string  `yaml:"name"`
	Props      []_Prop `yaml:"props,omitempty"`
	Embeds     []string
	Implements []string
}

type _Resource struct {
	Name       string       `yaml:"name"`
	Spec       string       `yaml:"spec,omitempty"`
	Status     string       `yaml:"status,omitempty"`
	ApiMethods []_ApiMethod `yaml:"apimethods,omitempty"`
}

func loadModel(path string) ([]_Struct, []_Resource) {
	yamls := yamlFiles(path)
	structs := []_Struct{}
	resources := []_Resource{}

	for _, y := range yamls {
		model, err := readModel(y)
		if err != nil {
			log.Fatal(err)
		}

		for _, m := range model.Types {
			if m.Kind == "Struct" {
				structs = append(structs, _Struct{
					Name:  m.Name,
					Props: capitalizeProps(m.Props),
				})
				continue
			}
			if m.Kind == "Resource" {
				resources = append(resources, _Resource{
					Name:       m.Name,
					Spec:       m.Spec,
					Status:     m.Status,
					ApiMethods: m.ApiMethods,
				})
				continue
			}
		}
	}

	return structs, resources
}

func readModel(path string) (*_Model, error) {
	log.Printf("reading model %s ", path)

	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := _Model{}
	err = yaml.Unmarshal(yfile, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func capitalizeProps(l []_Prop) []_Prop {
	res := []_Prop{}
	for _, p := range l {
		res = append(res,
			_Prop{
				Prop:    capitalize(p.Prop),
				Json:    decapitalize(p.Prop),
				Type:    p.Type,
				Default: p.Default,
			})
	}
	return res
}
