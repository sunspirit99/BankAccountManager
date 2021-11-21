package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "BankAccountManager/prototype"
)

func main() {
	fmt.Println("Hello I am a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	//CREATE ACCOUNT

	req := &pb.TransferRequest{
		Transaction: &pb.Transaction{
			From:   1,
			To:     2,
			Amount: 1000,
		},
	}

	// req1 := &pb.AccountListRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.Acc_Transfer(ctx, req)
	// res, err := c.Acc_List(ctx, req1)

	if err != nil {
		log.Fatalf("Error while creating user profile %v", err)
	}
	fmt.Println(res)
}
