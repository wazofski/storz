package utils

import (
	"encoding/json"
	"log"

	"github.com/wazofski/store"
)

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
