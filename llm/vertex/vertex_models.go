package vertex

func (client *Client) ProvideModel(name string) bool {
	switch name {
	case "gemini-pro":
		return true
	default:
		return false
	}
}
