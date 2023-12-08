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
	// migration.Storage(supa, app)

	//connection := os.Getenv("NEON_DB")
	////connection := os.Getenv("AWS_BRAINBOOST_DB")
	//con, err := pgxpool.New(ctx, connection)
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer con.Close()
	//
	//err = con.Ping(ctx)
	//if err != nil {
	//	log.Fatalf("did not ping: %v", err)
	//}
	//
	//log.Printf("connected")

	//migration.UpdateCollections(ctx, con)

	migration.PineconeImport(ctx)
}
