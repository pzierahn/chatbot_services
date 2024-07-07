package openai

import (
	_ "embed"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/sashabaranov/go-openai"
	"reflect"
	"testing"
)

//go:embed openai_transform_test_llm.json
var llmMessagesByt []byte

//go:embed openai_transform_test_openai.json
var openaiMessagesByt []byte

func Test_messagesToOpenAI(t *testing.T) {
	var messages []*llm.Message
	err := json.Unmarshal(llmMessagesByt, &messages)
	if err != nil {
		t.Fatal(err)
	}

	claudeMessages := messagesToOpenAI(messages)

	var expected []openai.ChatCompletionMessage
	err = json.Unmarshal(openaiMessagesByt, &expected)
	if err != nil {
		t.Fatal(err)
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(claudeMessages, expected) {
		t.Fatalf("Expected %v, got %v", expected, claudeMessages)
	}
}

func Test_openaiToMessages(t *testing.T) {
	var messages []openai.ChatCompletionMessage
	err := json.Unmarshal(openaiMessagesByt, &messages)
	if err != nil {
		t.Fatal(err)
	}

	llmMessages := openaiToMessages(messages)

	var expected []*llm.Message
	err = json.Unmarshal(llmMessagesByt, &expected)
	if err != nil {
		t.Fatal(err)
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(llmMessages, expected) {
		byt, _ := json.MarshalIndent(llmMessages, "", "  ")
		t.Log("Got:", string(byt))

		t.Fatalf("Expected %v, got %v", expected, llmMessages)
	}
}
