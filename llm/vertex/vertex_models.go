package vertex

const modelPrefix = "google."

func (client *Client) ProvidesModel(name string) bool {
	switch name {
	case modelPrefix + "gemini-pro":
		return true
	default:
		return false
	}
}
