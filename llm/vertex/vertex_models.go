package vertex

const modelPrefix = "google."

func (client *Client) ProvideModel(name string) bool {
	switch name {
	case modelPrefix + "gemini-pro":
		return true
	default:
		return false
	}
}
