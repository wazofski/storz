{{ $name := .Name }} 
type _{{ $name }} struct {
	{{ range .Embeds }}	{{ . }}
	{{ end }}
	{{ range .Props }}	{{ .Prop }}_ *{{.Type}} `json:"{{ .Json }}"`
	{{ end }}
}

{{ range .Props }}
func (entity *_{{$name}}) Set{{ .Prop }}(val {{.Type}}) {
	entity.{{.Prop}}_ = &val
}

func (entity *_{{$name}}) {{ .Prop }}() {{.Type}}{
	return *entity.{{.Prop}}_
}
{{ end }}

func {{ .Name }}Factory() {{.Name}} {
	{{ range .Props }}{{ .Prop }}_ := {{ .Default }}
	{{ end }}
	
	return &_{{ .Name }} {
		{{ range .Props }}{{ .Prop }}_: &{{ .Prop }}_,
		{{ end }}
	}
}

