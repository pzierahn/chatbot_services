package test

import "encoding/json"

func prettify(obj interface{}) string {
	byt, _ := json.MarshalIndent(obj, "", "  ")
	return string(byt)
}
