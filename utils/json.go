package utils

import (
	"encoding/json"
	"log"
	"os"
)

func Prettify(obj interface{}) string {
	byt, _ := json.MarshalIndent(obj, "", "  ")
	return string(byt)
}

func WriteJson(path string, obj interface{}) {
	byt, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(path, byt, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
