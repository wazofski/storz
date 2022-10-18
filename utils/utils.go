package utils

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/wazofski/storz/constants"
	"github.com/wazofski/storz/store"
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

func ReadStream(r io.ReadCloser) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

func UnmarshalObject(body []byte, schema store.SchemaHolder, kind string) (store.Object, error) {
	resource := schema.ObjectForKind(kind)
	err := json.Unmarshal(body, &resource)

	return resource, err
}

func ObjeectKind(response []byte) string {
	resource := _Resource{}
	err := json.Unmarshal(response, &resource)
	if err != nil {
		return ""
	}

	if resource.Metadata == nil {
		return ""
	}

	return resource.Metadata.(map[string]interface{})["kind"].(string)
}

func PP(obj store.Object) string {
	jsn, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	// log.Println(obj)

	return string(jsn)
}

func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func Serialize(mo store.Object) ([]byte, error) {
	if mo == nil {
		return nil, constants.ErrObjectNil
	}

	return json.Marshal(mo)
}
