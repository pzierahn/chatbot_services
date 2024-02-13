package bedrock

func (client *Client) ProvidesModel(name string) bool {
	switch name {
	case "amazon.titan-text-express-v1":
		return true
	case "anthropic.claude-v2":
		return true
	default:
		return false
	}
}
