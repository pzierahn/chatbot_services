package collections

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	storage_go "github.com/supabase-community/storage-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"sort"
	"testing"
)

func TestService_GetAll(t *testing.T) {
	setupTestService(t)
	defer teardownTestService(t)

	ctx := context.Background()

	collections := []*pb.Collection{
		{
			Name: "Collection 2",
		},
		{
			Name: "Collection 1",
		},
		{
			Name: "Collection 3",
		},
	}

	for _, collection := range collections {
		_, err := testService.Create(ctx, collection)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Sort by name after creation
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].Name < collections[j].Name
	})

	data, err := testService.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	if len(data.Items) != len(collections) {
		t.Fatal("wrong number of collections")
	}

	for i, collection := range collections {
		if collection.Name != data.Items[i].Name {
			t.Fatal("wrong collection name")
		}
	}
}

func TestService_GetAll1(t *testing.T) {
	setupTestService(t)
	defer teardownTestService(t)

	ctx := context.Background()

	// Create test users
	users := []struct {
		id          uuid.UUID
		collections []*pb.Collection
	}{
		{
			id: uuid.New(),
			collections: []*pb.Collection{
				{
					Name: "Collection 1",
				},
				{
					Name: "Collection 2",
				},
			},
		},
		{
			id: uuid.New(),
			collections: []*pb.Collection{
				{
					Name: "Data 1",
				},
				{
					Name: "Data 2",
				},
			},
		},
	}

	// Create users and collections
	for _, user := range users {
		service := &Service{
			auth: auth.WithUser(user.id),
			db:   testConn,
		}

		for _, collection := range user.collections {
			_, err := service.Create(ctx, collection)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	type fields struct {
		auth    auth.Service
		db      *pgxpool.Pool
		storage *storage_go.Client
	}

	tests := []struct {
		name    string
		fields  fields
		want    *pb.Collections
		wantErr bool
	}{
		{
			name: "GetAll user 1",
			fields: fields{
				auth: auth.WithUser(users[0].id),
				db:   testConn,
			},
			want: &pb.Collections{
				Items: []*pb.Collections_Collection{
					{
						Name: "Collection 1",
					},
					{
						Name: "Collection 2",
					},
				},
			},
		},
		{
			name: "GetAll user 2",
			fields: fields{
				auth: auth.WithUser(users[1].id),
				db:   testConn,
			},
			want: &pb.Collections{
				Items: []*pb.Collections_Collection{
					{
						Name: "Data 1",
					},
					{
						Name: "Data 2",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Service{
				auth: tt.fields.auth,
				db:   tt.fields.db,
			}
			got, err := server.GetAll(ctx, &emptypb.Empty{})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i, collection := range got.Items {
				if collection.Name != tt.want.Items[i].Name {
					t.Errorf("GetAll() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
