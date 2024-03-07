package anthropic

const prefix = "anthropic."

const (
	OPUS = "claude-3-opus-20240229"
)

func (client *Client) ProvidesModel(name string) bool {
	switch {
	case name == OPUS:
		return true
	case prefix+name == OPUS:
		return true
	default:
		return false
	}
}
