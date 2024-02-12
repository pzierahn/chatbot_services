package bedrock

import (
	"context"
	"fmt"
	"github.com/pzierahn/chatbot_services/llm"
	"strings"
)

func (client *Client) GenerateCompletion(ctx context.Context, req *llm.GenerateRequest) (*llm.GenerateResponse, error) {
	switch {
	case strings.HasPrefix(req.Model, "amazon"):
		return client.generateCompletionTitan(ctx, req)
	case strings.HasPrefix(req.Model, "anthropic"):
		return client.generateCompletionAnthropic(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported model: %s", req.Model)
	}
}
