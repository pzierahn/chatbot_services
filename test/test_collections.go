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

	setup.report.ExpectError("collection_create_without_auth", func() error {
		_, err := setup.collections.Create(context.Background(), &pb.Collection{
			Name: "Test Collection",
		})
		return err
	})

	setup.report.Run("collection_create", func() error {
		coll, err := setup.collections.Create(ctx, &pb.Collection{
			Name: "Test Collection",
		})
		if err != nil {
			return err
		}

		colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		if len(colls.Items) != 1 {
			return fmt.Errorf("collection not created")
		}

		if colls.Items[0].Id != coll.Id {
			return fmt.Errorf("collection id mismatch")
		}

		if colls.Items[0].Name != coll.Name {
			return fmt.Errorf("collection name mismatch")
		}

		return nil
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

	setup.report.ExpectError("collection_update_without_auth", func() error {
		_, err = setup.collections.Update(context.Background(), update)
		return err
	})

	setup.report.Run("collection_update_valid", func() error {
		_, err = setup.collections.Update(ctx, update)
		if err != nil {
			return err
		}

		colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		if len(colls.Items) != 1 {
			return fmt.Errorf("expected 1 collection")
		}

		if colls.Items[0].Id != coll.Id {
			return fmt.Errorf("collection id mismatch")
		}

		if colls.Items[0].Name != update.Name {
			return fmt.Errorf("collection name mismatch")
		}

		return nil
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

	colls, err := setup.collections.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range colls.Items {
		if c.Id != coll.Id {
			continue
		}

		log.Println("Collection found:", c.Name)

		if c.Name != coll.Name {
			log.Fatal("Collection name mismatch")
		}

		break
	}

	_, err = setup.collections.Delete(ctx, coll)
	if err != nil {
		log.Fatal(err)
	}

	colls, err = setup.collections.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range colls.Items {
		if c.Id == coll.Id {
			log.Fatalf("Collection %s not deleted", c.Name)
		}
	}
}
