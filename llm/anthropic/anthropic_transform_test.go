package anthropic

import (
	_ "embed"
	"encoding/json"
	"github.com/pzierahn/chatbot_services/llm"
	"reflect"
	"testing"
)

//go:embed anthropic_transform_test_llm.json
var llmMessagesByt []byte

//go:embed anthropic_transform_test_claude.json
var claudeMessagesByt []byte

func Test_transformToClaude(t *testing.T) {
	var messages []*llm.Message
	err := json.Unmarshal(llmMessagesByt, &messages)
	if err != nil {
		t.Fatal(err)
	}

	claudeMessages, err := transformToClaude(messages)
	if err != nil {
		t.Fatal(err)
	}

	var expected []ClaudeMessage
	err = json.Unmarshal(claudeMessagesByt, &expected)
	if err != nil {
		t.Fatal(err)
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(claudeMessages, expected) {
		t.Fatalf("Expected %v, got %v", expected, claudeMessages)
	}
}

func Test_claudeToMessages(t *testing.T) {
	var messages []ClaudeMessage
	err := json.Unmarshal(claudeMessagesByt, &messages)
	if err != nil {
		t.Fatal(err)
	}

	llmMessages, err := claudeToMessages(messages)
	if err != nil {
		t.Fatal(err)
	}

	var expected []*llm.Message
	err = json.Unmarshal(llmMessagesByt, &expected)
	if err != nil {
		t.Fatal(err)
	}

	// Deep equal comparison with reflect.DeepEqual
	if !reflect.DeepEqual(llmMessages, expected) {
		t.Fatalf("Expected %v, got %v", expected, llmMessages)
	}
}
