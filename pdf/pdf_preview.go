package pdf

import (
	"bytes"
	"context"
	"os/exec"
)

func GetPageAsImage(ctx context.Context, input []byte) (byt []byte, err error) {
	cmd := exec.CommandContext(ctx, "pdftoppm",
		"-jpeg",
		"-f", "1", "-l", "1",
		"-scale-to", "120",
		"-")
	cmd.Stdin = bytes.NewReader(input)

	return cmd.CombinedOutput()
}
