package main

import (
	"encoding/json"
	"github.com/pzierahn/brainboost/test"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	testing := test.NewTestSetup()
	defer testing.Close()

	testing.CollectionCreate()
	testing.CollectionRename()
	testing.CollectionDelete()

	testing.DocumentsIndex()
	testing.DocumentsList()
	testing.DocumentsDelete()
	testing.DocumentsSearch()
	testing.DocumentsUpdate()

	testing.ChatGenerate()
	testing.ChatHistory()

	testing.Account()

	byt, _ := json.MarshalIndent(testing.Report, "", "  ")
	log.Println(string(byt))
}
