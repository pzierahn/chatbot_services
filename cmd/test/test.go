package main

import (
	"github.com/pzierahn/brainboost/test"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	testing := test.NewTestSetup()
	defer testing.Close()

	testing.CollectionCreate()
	testing.CollectionRename()
}
