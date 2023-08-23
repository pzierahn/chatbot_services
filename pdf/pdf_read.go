package pdf

import (
	"context"
	"os/exec"
	"strings"
)

// ReadPages reads a PDF file and returns its contents as a slice of strings.
func ReadPages(ctx context.Context, filename string) (result []string, err error) {
	cmd := exec.CommandContext(ctx, "pdftotext",
		"-layout",
		"-htmlmeta",
		filename,
		"-",
	)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	text := string(data)

	return strings.Split(text, "\f"), nil
}
