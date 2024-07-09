package utils

import (
	"context"
	"io"
	"jaytaylor.com/html2text"
	"net/http"
)

func Scrape(ctx context.Context, url string) (text string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	htmlText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return html2text.FromString(string(htmlText), html2text.Options{
		PrettyTables: false,
		OmitLinks:    true,
	})
}
