package search

import "context"

type Fragment struct {
	Id           string `json:"id,omitempty"`
	Text         string `json:"text,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	DocumentId   string `json:"document_id,omitempty"`
	CollectionId string `json:"collection_id,omitempty"`
	Position     uint32 `json:"position,omitempty"`
}

type Query struct {
	UserId       string  `json:"user_id,omitempty"`
	CollectionId string  `json:"collection_id,omitempty"`
	Query        string  `json:"query,omitempty"`
	Limit        uint32  `json:"limit,omitempty"`
	Threshold    float32 `json:"threshold,omitempty"`
}

type Result struct {
	Id         string  `json:"id,omitempty"`
	Text       string  `json:"text,omitempty"`
	DocumentId string  `json:"document_id,omitempty"`
	Position   uint32  `json:"position,omitempty"`
	Score      float32 `json:"score,omitempty"`
}

type Index interface {
	Search(context.Context, Query) ([]*Result, error)
	Upsert(context.Context, []*Fragment) error
	DeleteCollection(ctx context.Context, userId, collectionId string) error
	DeleteDocument(ctx context.Context, userId, documentId string) error
	Close() error
}
