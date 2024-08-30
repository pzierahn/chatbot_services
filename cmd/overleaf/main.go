package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const texDirectory = "tmp_overleaf"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if _, err := os.Stat(texDirectory); os.IsNotExist(err) {
		// Checkout the overleaf project

		err = os.Mkdir(texDirectory, 0755)
		if err != nil {
			log.Fatal(err)
		}
		//defer func() { _ = os.RemoveAll("tmp_overleaf") }()

		token := os.Getenv("OVERLEAF_TOKEN")
		url := fmt.Sprintf("https://git:%s@git.overleaf.com/6617becf3bf104428e4f748c", token)
		log.Printf("url: %s\n", url)

		ctx := context.Background()
		cmd := exec.CommandContext(ctx, "git", "clone", url, ".")
		cmd.Dir = texDirectory

		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Pull the overleaf project
		ctx := context.Background()
		cmd := exec.CommandContext(ctx, "git", "pull")
		cmd.Dir = texDirectory
	}

	// Read all the tex files in the part directory
	files, err := os.ReadDir(filepath.Join(texDirectory, "parts"))
	if err != nil {
		log.Fatal(err)
	}

	var texParts []string
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".tex") {
			// Skip directories and non-tex files
			continue
		}

		log.Printf("file: %s\n", file.Name())
		// Read the content of the file
		path := filepath.Join(texDirectory, "parts", file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		texParts = append(texParts, string(content))
	}

	texContent := strings.Join(texParts, "\n\n")
	//log.Printf("tex content: %s", texContent)

	// Write the content to a new file
	err = os.WriteFile("main.tex", []byte(texContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
