package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pzierahn/chatbot_services/datastore"
	"time"
)

func main() {
	ctx := context.Background()
	store, err := datastore.New(ctx)
	if err != nil {
		panic(err)
	}

	err = store.InsertPayment(ctx, &datastore.Payments{
		Id:     uuid.New(),
		UserId: "",
		Amount: 0,
		Date:   time.Now(),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Payment inserted")
}
