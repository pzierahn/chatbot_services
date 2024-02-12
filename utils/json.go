package utils

import "encoding/json"

func Prettify(obj interface{}) string {
	byt, _ := json.MarshalIndent(obj, "", "  ")
	return string(byt)
}
