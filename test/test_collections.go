package test

import (
	"context"
	supa "github.com/nedpals/supabase-go"
	pb "github.com/pzierahn/brainboost/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func (setup *Setup) CollectionCreate() {

	user := setup.CreateUser()
	defer setup.DeleteUser(user.Id)

	supabase := supa.CreateClient(setup.SupabaseUrl, setup.Token)
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

	user := setup.CreateUser()
	defer setup.DeleteUser(user.Id)

	supabase := supa.CreateClient(setup.SupabaseUrl, setup.Token)
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
