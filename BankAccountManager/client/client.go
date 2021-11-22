package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	// CREATE ACCOUNT

	// req := &pb.TransferRequest{
	// 	Transaction: &pb.Transaction{
	// 		From:   1,
	// 		To:     2,
	// 		Amount: 1000,
	// 	},
	// }

	// req1 := &pb.AccountListRequest{}

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// res, err := c.Acc_Transfer(ctx, req)
	// res, err := c.Acc_List(ctx, req1)

	// if err != nil {
	// 	log.Fatalf("Error while creating account %v", err)
	// }
	// fmt.Println(res)

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
					Balance:     20000,
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
		fmt.Printf("Account No.%v Created : %v \n", i, reply.GetAccount())
	}

	// from := uint64(rand.Int63n(4) + 1)
	// to := uint64(rand.Int63n(4) + 1)
	// for i := 0; i < 1000; i++ {
	// 	req := &pb.TransferRequest{
	// 		Transaction: &pb.Transaction{
	// 			From:   from,
	// 			Amount: 1,
	// 			To:     to,
	// 		},
	// 	}

	// 	res, err := c.Acc_Transfer(context.Background(), req)
	// 	if err != nil {
	// 		fmt.Printf("Failed to transfer from [AcountID = %v] to [AcountID = %v]", from, to)
	// 	}
	// 	fmt.Println(res)
	// }

}
