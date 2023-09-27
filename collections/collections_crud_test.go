package collections

import (
	"context"
	pb "github.com/pzierahn/brainboost/proto"
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
