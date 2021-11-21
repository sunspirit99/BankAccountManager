package main

import (
	"context"
	"fmt"
	"io"
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

	stream, err := c.Acc_Create(context.Background())
	if err != nil {
		log.Fatal("Failed to request creating account to server :", err)
	}

	var count uint64 = 1
	go func() {
		for {
			req := &pb.AccountRequest{
				Account: &pb.Account{
					Id:          count,
					Name:        "Linh",
					Address:     "Tay Ho",
					PhoneNumber: "0123456789",
					Balance:     200000,
					Status:      1,
					Createtime:  timestamppb.Now(),
				},
			}
			count += 1
			if err := stream.Send(req); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Nanosecond)
		}
	}()

	for i := 0; i < 1000; i++ {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("Account No.%v Created  :\n", reply.GetAccount())
	}

}
