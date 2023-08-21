package pdf

import (
	"context"
	"fmt"
	"os/exec"
)

// ReadFile reads a PDF file and returns its contents as a byte slice.
func ReadFile(ctx context.Context, filename string) (result []byte, err error) {
	cmd := exec.CommandContext(ctx, "pdftotext", filename, "-")
	return cmd.CombinedOutput()
}

// ReadFileString reads a PDF file and returns its contents as a string.
func ReadFileString(ctx context.Context, filename string) (result string, err error) {
	data, err := ReadFile(ctx, filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadPages reads a PDF file and returns its contents as a slice of strings.
func ReadPages(ctx context.Context, filename string, pages int) (result []string, err error) {
	for page := 1; page <= pages; page++ {

		cmd := exec.CommandContext(ctx, "pdftotext",
			"-f", fmt.Sprint(page),
			"-l", fmt.Sprint(page),
			"-layout",
			filename,
			"-",
		)
		data, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		result = append(result, string(data))
	}

	return result, nil
}
