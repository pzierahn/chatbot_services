package test

import (
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (setup *Setup) CollectionCreate() {

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

	coll.Name = "Test Collection 2"
	_, err = setup.collections.Update(ctx, coll)
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
