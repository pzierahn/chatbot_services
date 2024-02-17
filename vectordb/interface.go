package vectordb

type Vector struct {
	Id           string
	DocumentId   string
	UserId       string
	CollectionId string
	Text         string
	Vector       []float32
	Score        float32
}

type SearchQuery struct {
	UserId       string
	CollectionId string
	Vector       []float32
	Limit        int
	Threshold    float32
}

type DB interface {
	Close() error
	Delete(ids []string) error
	Upsert(items []*Vector) error
	Search(query SearchQuery) ([]*Vector, error)
}
