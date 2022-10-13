package utils

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/wazofski/store"
)

type _Resource struct {
	Metadata interface{} `json:"metadata,omitempty"`
	Spec     interface{} `json:"spec,omitempty"`
	Status   interface{} `json:"status,omitempty"`
}

func CloneObject(obj store.Object, schema store.SchemaHolder) store.Object {
	ret := schema.ObjectForKind(obj.Metadata().Kind())
	jsn, err := json.Marshal(obj)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(jsn, &ret)
	if err != nil {
		log.Panic(err)
	}
	return ret
}

func ReadStream(body io.ReadCloser) ([]byte, error) {
	var b strings.Builder
	var n int
	var err error
	data := make([]byte, 128)
	for n, err = body.Read(data); err == nil; {
		b.WriteString(string(data[:n]))
	}

	return []byte(b.String()), err
}

func UnmarshalObject(body []byte, schema store.SchemaHolder) (store.Object, error) {
	resource := schema.ObjectForKind(ObjeectKind(body))
	err := json.Unmarshal(body, &resource)

	return resource, err
}

func ObjeectKind(response []byte) string {
	resource := _Resource{}
	err := json.Unmarshal(response, &resource)
	if err != nil {
		log.Printf("Error parsing %s: %s", string(response), err)
		return ""
	}

	return resource.Metadata.(map[string]interface{})["kind"].(string)
}
