package migration

import (
	"github.com/pzierahn/chatbot_services/search"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migrator struct {
	Database *mongo.Client
	Search   search.Index
}
