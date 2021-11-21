package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/grpc"

	pb "BankAccountManager/prototype"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateUserProfile(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	//CREATE ACCOUNT

	req := &pb.AccountRequest{
		Account: &pb.Account{
			Id:          5,
			Name:        "Bin",
			Address:     "Tay Ho",
			PhoneNumber: "0979667408",
			Balance:     20000,
			Status:      0,
			Createtime:  timestamppb.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.Acc_Create(ctx, req)

	if err != nil {
		log.Fatalf("Error while creating account %v", err)
	}
	fmt.Println(res)
}
