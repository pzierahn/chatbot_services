package qdrant

import (
	"crypto/tls"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/llm/voyageai"
	"github.com/pzierahn/chatbot_services/vectordb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type DB struct {
	conn      *grpc.ClientConn
	apiKey    string
	namespace string
	embedding llm.Embedding
	dimension int
	queue     chan *vectordb.Fragment
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func New() (*DB, error) {
	apiKey := os.Getenv("CHATBOT_QDRANT_KEY")
	target := os.Getenv("CHATBOT_QDRANT_URL")

	var opts []grpc.DialOption
	if os.Getenv("CHATBOT_QDRANT_INSECURE") == "true" {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	voyage, err := voyageai.New(voyageai.ModelVoyageLarge2)
	if err != nil {
		return nil, err
	}

	client := &DB{
		conn:      conn,
		apiKey:    apiKey,
		namespace: "documents_v2",
		dimension: voyageai.DimensionVoyageLarge2,
		embedding: voyage,
	}

	err = client.Init()
	if err != nil {
		return nil, err
	}

	return client, nil
}
