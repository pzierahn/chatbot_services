package pdf

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

func findPageNumber(data []byte) (pages int, err error) {
	// Define the regex pattern to match the "Pages" line and its value
	pagesPattern := regexp.MustCompile(`Pages:\s+(\d+)`)

	// Find all matches in the text
	matches := pagesPattern.FindSubmatch(data)

	if len(matches) >= 2 {
		return strconv.Atoi(string(matches[1]))
	} else {
		return 0, fmt.Errorf("pages not found")
	}
}

// GetPageCount returns the number of pages in a PDF file.
func GetPageCount(ctx context.Context, filename string) (pages int, err error) {
	cmd := exec.CommandContext(ctx, "pdfinfo", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	return findPageNumber(out)
}
