package pinecone_search

import (
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/pzierahn/chatbot_services/llm"
	"github.com/pzierahn/chatbot_services/search"
	"log"
	"os"
)

type Search struct {
	conn          *pinecone.Client
	namespace     string
	embedding     llm.Embedding
	fastEmbedding *search.ParallelEmbedding
	dimension     int
}

func New(engine llm.Embedding, namespace string) (*Search, error) {
	clientParams := pinecone.NewClientParams{
		ApiKey: os.Getenv("PINECONE_API_KEY"),
	}

	pc, err := pinecone.NewClient(clientParams)

	if err != nil {
		log.Fatalf("Failed to create Client: %v", err)
	}

	fastEmbedding := search.NewParallelEmbedding(engine, 10, 100)

	client := &Search{
		conn:          pc,
		namespace:     namespace,
		embedding:     engine,
		dimension:     engine.GetEmbeddingDimension(),
		fastEmbedding: fastEmbedding,
	}

	err = client.Init()
	if err != nil {
		return nil, err
	}

	return client, nil
}
