package test

import (
	"context"
	"fmt"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (setup *Setup) CollectionCreate() {

	ctx, userId := setup.createRandomSignIn()
	defer setup.DeleteUser(userId)

	setup.Report.Run("collection_create_without_auth", func(t testing) bool {
		_, err := setup.collections.Create(context.Background(), &pb.Collection{
			Name: "Test Collection",
		})
		return t.expectError(err)
	})

	setup.Report.Run("collection_create", func(t testing) bool {
		coll, err := setup.collections.Create(ctx, &pb.Collection{
			Name: "Test Collection",
		})
		if err != nil {
			return t.fail(err)
		}

		colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
		if err != nil {
			return t.fail(err)
		}

		if len(colls.Items) != 1 {
			return t.fail(fmt.Errorf("collection not created"))
		}

		if colls.Items[0].Id != coll.Id {
			return t.fail(fmt.Errorf("collection id mismatch"))
		}

		if colls.Items[0].Name != coll.Name {
			return t.fail(fmt.Errorf("collection name mismatch"))
		}

		return t.pass()
	})
}

func (setup *Setup) CollectionRename() {

	ctx, userId := setup.createRandomSignIn()
	defer setup.DeleteUser(userId)

	coll, err := setup.collections.Create(ctx, &pb.Collection{
		Name: "Test Collection",
	})
	if err != nil {
		log.Fatal(err)
	}

	update := &pb.Collection{
		Id:   coll.Id,
		Name: "Test Collection 2",
	}

	setup.Report.Run("collection_update_without_auth", func(t testing) bool {
		_, err = setup.collections.Update(context.Background(), update)
		return t.expectError(err)
	})

	setup.Report.Run("collection_update_valid", func(t testing) bool {
		_, err = setup.collections.Update(ctx, update)
		if err != nil {
			return t.fail(err)
		}

		colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
		if err != nil {
			return t.fail(err)
		}

		if len(colls.Items) != 1 {
			return t.fail(fmt.Errorf("expected 1 collection"))
		}

		if colls.Items[0].Id != coll.Id {
			return t.fail(fmt.Errorf("collection id mismatch"))
		}

		if colls.Items[0].Name != update.Name {
			return t.fail(fmt.Errorf("collection name mismatch"))
		}

		return t.pass()
	})
}

func (setup *Setup) CollectionDelete() {

	ctx, userId := setup.createRandomSignIn()
	defer setup.DeleteUser(userId)

	coll, err := setup.collections.Create(ctx, &pb.Collection{
		Name: "Test Collection",
	})
	if err != nil {
		log.Fatal(err)
	}

	setup.Report.Run("collection_delete_without_auth", func(t testing) bool {
		_, err = setup.collections.Delete(context.Background(), coll)
		return t.expectError(err)
	})

	setup.Report.Run("collection_delete_invalid", func(t testing) bool {
		_, err = setup.collections.Delete(ctx, &pb.Collection{})
		return t.expectError(err)
	})

	setup.Report.Run("collection_delete_valid", func(t testing) bool {
		_, err = setup.collections.Delete(ctx, coll)
		if err != nil {
			return t.fail(err)
		}

		colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
		if err != nil {
			return t.fail(err)
		}

		for _, c := range colls.Items {
			if c.Id == coll.Id {
				return t.fail(fmt.Errorf("collection %s not deleted", c.Name))
			}
		}

		return t.pass()
	})
}
