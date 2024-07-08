package vectordb

type Fragment struct {
	Id           string `json:"id,omitempty"`
	DocumentId   string `json:"document_id,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	CollectionId string `json:"collection_id,omitempty"`
	Text         string `json:"text,omitempty"`
}

type SearchQuery struct {
	UserId       string  `json:"user_id,omitempty"`
	CollectionId string  `json:"collection_id,omitempty"`
	Query        string  `json:"query,omitempty"`
	Limit        int     `json:"limit,omitempty"`
	Threshold    float32 `json:"threshold,omitempty"`
}

type SearchResults struct {
	Fragments []*Fragment `json:"fragments,omitempty"`
	Scores    []float32   `json:"scores,omitempty"`
}

type DB interface {
	Search(query SearchQuery) (*SearchResults, error)
	Upsert(items []*Fragment) error
	Delete(ids []string) error
	Close() error
}
