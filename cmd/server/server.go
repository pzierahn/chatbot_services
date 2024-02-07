package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/chatbot_services/account"
	"github.com/pzierahn/chatbot_services/auth"
	"github.com/pzierahn/chatbot_services/chat"
	"github.com/pzierahn/chatbot_services/collections"
	"github.com/pzierahn/chatbot_services/documents"
	"github.com/pzierahn/chatbot_services/llm/openai"
	"github.com/pzierahn/chatbot_services/llm/vertex"
	pb "github.com/pzierahn/chatbot_services/proto"
	"github.com/pzierahn/chatbot_services/setup"
	"github.com/pzierahn/chatbot_services/vectordb/qdrant"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

const credentialsFile = "service_account.json"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	ctx := context.Background()

	var opts []option.ClientOption

	if _, err := os.Stat(credentialsFile); err == nil {
		serviceAccount := option.WithCredentialsFile(credentialsFile)
		opts = append(opts, serviceAccount)
	}

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		log.Fatalf("failed to create firebase app: %v", err)
	}

	firebaseStorage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("failed to create firebase storage client: %v", err)
	}

	bucket, err := firebaseStorage.Bucket("brainboost-399710.appspot.com")
	if err != nil {
		log.Fatalf("did not get bucket: %v", err)
	}

	addr := os.Getenv("CHATBOT_DB")
	db, err := pgxpool.New(ctx, addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer db.Close()

	err = setup.SetupTables(ctx, db)
	if err != nil {
		log.Fatalf("failed to setup tables: %v", err)
	}

	openaiService := openai.New()
	vertexService, err := vertex.New(ctx)
	if err != nil {
		log.Fatalf("failed to create vertex service: %v", err)
	}

	authService, err := auth.WithFirebase(ctx, app)
	if err != nil {
		log.Fatalf("failed to create auth service: %v", err)
	}

	//vecDB, err := pinecone.New()
	vecDB, err := qdrant.New()
	if err != nil {
		log.Fatalf("failed to create vector db: %v", err)
	}
	defer func() { _ = vecDB.Close() }()

	grpcServer := grpc.NewServer()
	collectionServer := collections.NewServer(authService, db, bucket, vecDB)
	pb.RegisterCollectionServiceServer(grpcServer, collectionServer)

	accountService := account.FromConfig(&account.Config{
		Auth: authService,
		DB:   db,
	})
	pb.RegisterAccountServiceServer(grpcServer, accountService)

	docsService := documents.FromConfig(&documents.Config{
		Auth:       authService,
		Account:    accountService,
		DB:         db,
		Embeddings: openaiService,
		Storage:    bucket,
		VectorDB:   vecDB,
	})
	pb.RegisterDocumentServiceServer(grpcServer, docsService)

	chatServer := chat.FromConfig(&chat.Config{
		DB:              db,
		Openai:          openaiService,
		Vertex:          vertexService,
		DocumentService: docsService,
		AccountService:  accountService,
		AuthService:     authService,
	})
	pb.RegisterChatServiceServer(grpcServer, chatServer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9055"
		log.Printf("defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting server on %v", listener.Addr().String())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
