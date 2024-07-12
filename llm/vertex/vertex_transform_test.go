package vertex

import (
	"cloud.google.com/go/vertexai/genai"
	_ "embed"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
	"reflect"
	"testing"
)

//go:embed vertex_transform_test_llm.json
var llmMessagesByt []byte

//go:embed vertex_transform_test_history.json
var vertexHistoryByt []byte

type fakeGenaiContent struct {
	Role  string
	Parts []interface{}
}

func (f *fakeGenaiContent) defuck() *genai.Content {
	var parts []genai.Part

	for _, part := range f.Parts {
		if str, ok := part.(string); ok {
			parts = append(parts, genai.Text(str))
		}

		obj, ok := part.(map[string]interface{})
		if !ok {
			continue
		}

		if args, isFuncCall := obj["Args"]; isFuncCall {
			name := obj["Name"].(string)
			argsx := args.(map[string]any)

			parts = append(parts, genai.FunctionCall{
				Name: name,
				Args: argsx,
			})
		}

		if response, isFuncResp := obj["Response"]; isFuncResp {
			name := obj["Name"].(string)
			responsex := response.(map[string]any)

			parts = append(parts, genai.FunctionResponse{
				Name:     name,
				Response: responsex,
			})
		}
	}

	return &genai.Content{
		Role:  f.Role,
		Parts: parts,
	}
}

func Test_transformToHistory(t *testing.T) {
	var messages []*llm.Message
	err := json.Unmarshal(llmMessagesByt, &messages)
	if err != nil {
		t.Fatal(err)
	}

	history, err := transformToHistory(messages)
	if err != nil {
		t.Fatal(err)
	}

	var defuckContent []fakeGenaiContent
	err = json.Unmarshal(vertexHistoryByt, &defuckContent)
	if err != nil {
		t.Fatal(err)
	}

	expected := make([]*genai.Content, len(defuckContent))
	for inx, fuck := range defuckContent {
		expected[inx] = fuck.defuck()
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(history, expected) {
		t.Fatalf("Expected %v, got %v", expected, history)
	}
}

func Test_transformToMessages(t *testing.T) {
	var defuckContent []fakeGenaiContent
	err := json.Unmarshal(vertexHistoryByt, &defuckContent)
	if err != nil {
		t.Fatal(err)
	}
	history := make([]*genai.Content, len(defuckContent))
	for inx, fuck := range defuckContent {
		history[inx] = fuck.defuck()
	}

	messages, err := transformToMessages(history)
	if err != nil {
		t.Fatal(err)
	}

	var expected []*llm.Message
	err = json.Unmarshal(llmMessagesByt, &expected)
	if err != nil {
		t.Fatal(err)
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(messages, expected) {
		bytex, _ := json.MarshalIndent(expected, "", "  ")
		bythi, _ := json.MarshalIndent(messages, "", "  ")

		t.Log("expected", string(bytex))
		t.Log("messages", string(bythi))

		t.Fatalf("Expected %v, got %v", expected, history)
	}
}
