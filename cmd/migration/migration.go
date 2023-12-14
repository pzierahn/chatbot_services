package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/pzierahn/brainboost/migration"
	"google.golang.org/api/option"
	"log"
)

var app *firebase.App

func init() {
	opt := option.WithCredentialsFile("service_account.json")

	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()

	//supa := migration.InitSupabase(ctx)
	//migration.Storage(supa, app)

	//connection := os.Getenv("BRAINBOOST_COCKROACH_DB")
	////connection := os.Getenv("NEON_DB")
	////connection := os.Getenv("AWS_BRAINBOOST_DB")
	//con, err := pgxpool.New(ctx, connection)
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer con.Close()

	//migration.MigratePayments(supa.DB, con)
	//migration.MigrateOpenaiUsage(supa.DB, con)
	//
	//migration.MigrateCollections(supa.DB, con)
	//migration.MigrateChatMessages(supa.DB, con)
	//migration.MigrateDocuments(supa.DB, con)
	//migration.MigrateDocumentsChunks(supa.DB, con)
	//migration.MigrateChatSources(supa.DB, con)

	migration.PineconeImport(ctx)
}
