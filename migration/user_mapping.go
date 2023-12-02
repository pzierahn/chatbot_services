package migration

import (
	"encoding/json"
	"log"
	"os"
)

func GetUserIdMapping() map[string]string {
	byt, err := os.ReadFile("user_mappings.json")
	if err != nil {
		log.Fatalf("Failed to read user_mappings.json: %v", err)
	}

	var mapping map[string]string
	err = json.Unmarshal(byt, &mapping)
	if err != nil {
		log.Fatalf("Failed to unmarshal user_mappings.json: %v", err)
	}

	return mapping
}
