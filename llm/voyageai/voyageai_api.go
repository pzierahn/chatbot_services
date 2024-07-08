package voyageai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	url = "https://api.voyageai.com/v1/embeddings"
)

const (
	InputTypeDocument = "document"
)

type Data struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type Usage struct {
	TotalTokens uint32 `json:"total_tokens"`
}

// Response is the response structure for the Voyage AI API.
type Response struct {
	Object string `json:"object"`
	Data   []Data `json:"data"`
	Model  string `json:"model"`
	Usage  Usage  `json:"usage"`
}

// Request is the request structure for the Voyage AI API.
type Request struct {
	Input     []string `json:"input"`
	Model     string   `json:"model"`
	InputType string   `json:"input_type"`
}

// callAPI calls the Voyage AI API with the given request and returns the response.
func (voyage *Client) callAPI(ctx context.Context, request *Request) (*Response, error) {

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+voyage.apiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data Response
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
