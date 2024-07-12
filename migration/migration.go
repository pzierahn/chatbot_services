package migration

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/datastore"
)

type Migrator struct {
	Legacy *pgxpool.Pool
	Next   *datastore.Service
}
