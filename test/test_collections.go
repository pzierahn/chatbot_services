package test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	supa "github.com/nedpals/supabase-go"
	"github.com/pzierahn/brainboost/auth"
	"github.com/pzierahn/brainboost/collections"
	pb "github.com/pzierahn/brainboost/proto"
	"github.com/pzierahn/brainboost/setup"
	storagego "github.com/supabase-community/storage-go"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"time"
)

func (service Service) getCollectionService() *collections.Service {

	storage := storagego.NewClient(
		os.Getenv("API_EXTERNAL_URL")+"/storage/v1",
		os.Getenv("SERVICE_ROLE_KEY"),
		nil)

	ctx := context.Background()
	db, err := pgxpool.New(ctx, "postgres://postgres:your-super-secret-and-long-postgres-password@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer db.Close()

	// Query time from the database
	var currentTime time.Time
	err = db.QueryRow(ctx, "select now()").Scan(&currentTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Current time:", currentTime)

	err = setup.SetupTables(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	supabaseAuth := auth.WithSupabase()
	return collections.NewServer(supabaseAuth, db, storage)
}

func (service Service) CreateCollection() {
	collectionService := service.getCollectionService()

	user := service.CreateUser()

	supabase := supa.CreateClient(service.SupabaseUrl, service.Token)
	details, err := supabase.Auth.SignIn(context.Background(), supa.UserCredentials{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"Authorization": []string{"Bearer " + details.AccessToken},
	})

	coll, err := collectionService.Create(ctx, &pb.Collection{
		Name: "Test Collection",
	})
	if err != nil {
		log.Fatal(err)
	}

	colls, err := collectionService.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range colls.Items {
		if c.Id == coll.Id {
			log.Println("Collection found:", c.Name)
		}
	}
}
