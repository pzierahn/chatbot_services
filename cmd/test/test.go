package main

import (
	"github.com/pzierahn/brainboost/test"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	service := test.Service{
		SupabaseUrl: os.Getenv("API_EXTERNAL_URL"),
		Token:       os.Getenv("SERVICE_ROLE_KEY"),
	}

	service.CreateCollection()
}
