package pdf

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// ReadPages reads a PDF file and returns its contents as a slice of strings.
func ReadPages(ctx context.Context, filename string) (result []string, err error) {
	cmd := exec.CommandContext(ctx, "pdftotext",
		//"-layout",
		//"-x", "0", "-y", "0", "-H", "500", "-W", "1000",
		filename,
		"-",
	)
	data, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	text := string(data)
	text = strings.TrimSpace(text)

	return strings.Split(text, "\f"), nil
}

func GetPagesFromBytes(ctx context.Context, data []byte) (result []string, err error) {
	cmd := exec.CommandContext(ctx, "pdftotext", "-", "-")
	cmd.Stdin = bytes.NewReader(data)

	pages, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	text := string(pages)
	text = strings.TrimSpace(text)

	return strings.Split(text, "\f"), nil
}
