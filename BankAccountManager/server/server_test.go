package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"

	"testing"

	"google.golang.org/grpc"

	pb "BankAccountManager/prototype"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var amount float32 = 1

func TestTransferManyToOne(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	for i := 0; i < 1000; i++ {
		var to uint64 = 1
		from := uint64(rand.Int63n(100) + 1)
		req := &pb.TransferRequest{
			Transaction: &pb.Transaction{
				From:   from,
				Amount: amount,
				To:     to,
			},
		}

		res, err := c.Acc_Transfer(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to transfer from [AcountID = %v] to [AcountID = %v]", from, to)
		}
		fmt.Println(res)
		log.Printf("Transaction No.%v \n", i)
	}
}

func TestTransferOneToMany(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	for i := 0; i < 1000; i++ {
		var from uint64 = 1
		to := uint64(rand.Int63n(100) + 1)
		req := &pb.TransferRequest{
			Transaction: &pb.Transaction{
				From:   from,
				Amount: amount,
				To:     to,
			},
		}

		res, err := c.Acc_Transfer(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to transfer from [AcountID = %v] to [AcountID = %v]", from, to)
		}
		fmt.Println(res)
		log.Printf("Transaction No.%v \n", i)
	}
}

func TestTransferOneToOne(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	for i := 0; i < 1000; i++ {
		from := uint64(rand.Int63n(4) + 1)
		to := uint64(rand.Int63n(5) + 5)
		req := &pb.TransferRequest{
			Transaction: &pb.Transaction{
				From:   from,
				Amount: amount,
				To:     to,
			},
		}

		res, err := c.Acc_Transfer(context.Background(), req)
		if err != nil {
			fmt.Printf("Failed to transfer from [AcountID = %v] to [AcountID = %v]", from, to)
		}
		fmt.Println(res)
		log.Printf("Transaction No.%v \n", i)
	}
}
func TestUpdateBalance(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	for i := 0; i < 1000; i++ {
		req := &pb.DepositRequest{
			Transaction: &pb.Transaction{
				From:   1,
				Amount: amount,
				To:     2,
			},
		}

		res, err := c.Acc_Deposit(context.Background(), req)
		if err != nil {
			fmt.Println("Failed to deposit to [AcountID = 1]")
		}
		fmt.Println(res)
	}
}
func TestCreateAccount(t *testing.T) {
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

	var count uint64 = 1003
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
		fmt.Printf("Account No.%v Created\n", reply.GetAccount())
	}

}
