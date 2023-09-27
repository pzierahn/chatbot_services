package collections

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pzierahn/brainboost/auth"
	pb "github.com/pzierahn/brainboost/proto"
	storage_go "github.com/supabase-community/storage-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"reflect"
	"testing"
)

func TestService_Create(t *testing.T) {
	setupTestService(t)
	defer teardownTestService(t)

	ctx := context.Background()

	tests := []struct {
		name       string
		collection *pb.Collection
		want       *pb.Collection
		wantErr    bool
	}{
		{
			name: "Create",
			collection: &pb.Collection{
				Name: "test",
			},
			want: &pb.Collection{
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "Create with bogus ID",
			collection: &pb.Collection{
				Id:   "asdfasdfasdf",
				Name: "test",
			},
			want: &pb.Collection{
				Name: "test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testService.Create(ctx, tt.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	setupTestService(t)
	defer teardownTestService(t)

	ctx := context.Background()

	testCollection, err := testService.Create(ctx, &pb.Collection{Name: "test"})
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}

	tests := []struct {
		name       string
		collection *pb.Collection
		want       *pb.Collection
		wantErr    bool
	}{
		{
			name: "Update without ID",
			collection: &pb.Collection{
				Name: "test",
			},
			wantErr: true,
		},
		{
			name: "Update with ID",
			collection: &pb.Collection{
				Id:   testCollection.Id,
				Name: "New Name",
			},
			want: &pb.Collection{
				Id:   testCollection.Id,
				Name: "New Name",
			},
			wantErr: false,
		},
		{
			name:       "Update nothing",
			collection: testCollection,
			want:       testCollection,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testService.Update(ctx, tt.collection)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_DeleteCollection(t *testing.T) {
	setupTestService(t)
	defer teardownTestService(t)

	ctx := context.Background()

	// Create test users
	collections := setupCollections(t)

	deletes := make(map[uuid.UUID]string)
	want := make(map[uuid.UUID]*pb.Collections)

	// Create users and collections
	for _, collection := range collections {
		service := &Service{
			auth: auth.WithUser(collection.userId),
			db:   testConn,
		}

		// Random delete collection
		index := rand.Intn(len(collection.collections))
		deleteCollection := collection.collections[index]
		deletes[collection.userId] = deleteCollection.Id

		_, err := service.Delete(ctx, deleteCollection)
		if err != nil {
			t.Errorf("Delete() error = %v", err)
			return
		}

		// Create expected result
		want[collection.userId] = &pb.Collections{}

		for _, coll := range collection.collections {
			if coll.Id != deleteCollection.Id {
				want[collection.userId].Items = append(want[collection.userId].Items, &pb.Collections_Collection{
					Id:   coll.Id,
					Name: coll.Name,
				})
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
		delete  uuid.UUID
		want    *pb.Collections
		wantErr bool
	}{
		{
			name: "GetAll user 1",
			fields: fields{
				auth: auth.WithUser(collections[0].userId),
				db:   testConn,
			},
			want: want[collections[0].userId],
		},
		{
			name: "GetAll user 2",
			fields: fields{
				auth: auth.WithUser(collections[1].userId),
				db:   testConn,
			},
			want: want[collections[0].userId],
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
