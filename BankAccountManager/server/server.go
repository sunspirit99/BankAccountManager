package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	c "BankAccountManager/config"
	"BankAccountManager/model"
	pb "BankAccountManager/prototype"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var conn *gorm.DB

const (
	port = ":50051"
)

type server struct {
	connector *gorm.DB
	pb.UnimplementedAccountServiceServer
}

func initDB() *gorm.DB {
	config :=
		c.Config{
			ServerName: "localhost:3306",
			User:       "root",
			Password:   "Password@99",
			DB:         "BankAccountManager",
		}

	connectionString := c.GetConnectionString(config)
	conn, err := Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}
	Migrate(&model.AccountORM{})
	return conn
}

// Connect creates MySQL connection
func Connect(connectionString string) (*gorm.DB, error) {
	var err error
	// dsn := "root:Password@99@tcp(127.0.0.1:3306)/BankAccountManager?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("DB Connection was successful!!")
	return conn, nil
}

// Migrate create/updates database table
func Migrate(table *model.AccountORM) {
	fmt.Println("Migrating DB ....")
	err := conn.AutoMigrate(&table)
	fmt.Println(err)
	log.Println("Table migrated")
}

func toORM(account *pb.Account) model.AccountORM {
	t := account.Createtime.AsTime()
	return model.AccountORM{
		Id:          account.Id,
		Name:        account.Name,
		Address:     account.Address,
		PhoneNumber: account.PhoneNumber,
		Balance:     account.Balance,
		Status:      int(account.Status),
		Createtime:  &t,
	}
}

func toPB(account *model.AccountORM) *pb.Account {
	t := timestamppb.New(*account.Createtime)
	return &pb.Account{
		Id:          account.Id,
		Name:        account.Name,
		Address:     account.Address,
		PhoneNumber: account.PhoneNumber,
		Balance:     account.Balance,
		Status:      pb.Account_STATE(account.Status),
		Createtime:  t,
	}
}

func (s *server) Acc_Create(stream pb.AccountService_Acc_CreateServer) error {
	db := s.connector
	// ORM := toORM(req.GetAccount())
	// error := db.Create(ORM).Error // Create a record in database
	// if error != nil {
	// 	fmt.Println("Query Error !")
	// 	return &pb.AccountResponse{
	// 		Message: "Failed",
	// 	}, error
	// }

	// res := &pb.AccountResponse{
	// 	Message: "Created",
	// }
	// return res, nil

	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("No more data !")
				return nil
			}
			return err
		}

		response := args.GetAccount()

		error := db.Create(toORM(response)).Error // Create a record in database
		if error != nil {
			fmt.Println("Query Error !")
			return error
		}

		res := &pb.AccountResponse{
			Message: "Created",
			Account: args.GetAccount(),
		}

		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func (s *server) Acc_Info(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	db := s.connector
	// ORM := toORM(account.GetAccount())
	id := req.GetAccount().GetId()
	var account *model.AccountORM
	err := db.First(&account, id).Error
	if err != nil {
		log.Printf("Account doesn't exist")
		return &pb.AccountResponse{}, err
	}

	res := &pb.AccountResponse{
		Message: "OK",
		Account: toPB(account),
	}
	return res, nil
}

func (s *server) Acc_List(ctx context.Context, req *pb.AccountListRequest) (*pb.AccountListResponse, error) {
	db := s.connector
	// ORM := toORM(account.GetAccount())

	var accountORM []model.AccountORM
	err := db.Find(&accountORM).Error
	if err != nil {
		log.Printf("Error when finding accounts !")
		return &pb.AccountListResponse{
			Message:  "Failed",
			Accounts: nil,
		}, err
	}

	var accounts []*pb.Account

	for _, acc := range accountORM {
		a := toPB(&acc)
		accounts = append(accounts, a)
	}

	return &pb.AccountListResponse{
		Message:  "OK",
		Accounts: accounts,
	}, nil

}

func (s *server) Acc_Update(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	db := s.connector
	// ORM := toORM(account.GetAccount())
	id := req.GetAccount().GetId()
	var account *model.AccountORM
	err := db.First(&account, id).Error
	if err != nil {
		log.Println("Account doesn't exist")
		return &pb.AccountResponse{}, err
	}

	new := req.GetAccount()
	if err := db.Updates(toORM(new)).Error; err != nil {
		log.Println("Failed to update account")
		return &pb.AccountResponse{
			Message: "Failed",
			Account: toPB(account),
		}, err
	}

	res := &pb.AccountResponse{
		Message: "OK",
		Account: new,
	}
	return res, nil
}

func (s *server) Acc_Delete(ctx context.Context, req *pb.AccountRequest) (*pb.AccountResponse, error) {
	db := s.connector

	id := req.GetAccount().GetId()
	var account *model.AccountORM
	err := db.First(&account, id).Error
	if err != nil {
		log.Printf("Account doesn't exist")
		return &pb.AccountResponse{}, err
	}

	err = db.Delete(account).Error
	if err != nil {
		log.Println("Failed to delete account")
		return &pb.AccountResponse{
			Message: "Failed",
			Account: req.GetAccount(),
		}, err
	}

	res := &pb.AccountResponse{
		Message: "OK",
	}
	return res, nil
}

func (s *server) Acc_Withdraw(ctx context.Context, req *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	db := s.connector

	id := req.GetTransaction().GetFrom()
	var account *model.AccountORM

	if req.GetTransaction().GetAmount() < 0 {
		log.Printf("Please enter an amount greater than 0")
		return &pb.WithdrawResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Account:     toPB(account),
		}, nil
	}

	err := db.First(&account, id).Error
	if err != nil {
		log.Printf("Account doesn't exist")
		return &pb.WithdrawResponse{
			Transaction: req.GetTransaction(),
		}, err
	}

	amount := account.Balance - req.GetTransaction().GetAmount()
	if amount < 0 {
		log.Printf("the balance is not enough to make this transaction")
		return &pb.WithdrawResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Account:     toPB(account),
		}, nil
	}

	account.Balance = amount
	if err := db.Updates(account).Error; err != nil {
		log.Println("Failed to withdraw")
		return &pb.WithdrawResponse{
			Message: "Failed",
			Account: toPB(account),
		}, err
	}

	res := &pb.WithdrawResponse{
		Message:     "OK",
		Transaction: req.GetTransaction(),
		Account:     toPB(account),
	}
	return res, nil
}

func (s *server) Acc_Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	db := s.connector

	id := req.GetTransaction().GetFrom()
	var account *model.AccountORM

	if req.GetTransaction().GetAmount() < 0 {
		log.Printf("Please enter an amount greater than 0")
		return &pb.DepositResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Account:     toPB(account),
		}, nil
	}

	err := db.First(&account, id).Error
	if err != nil {
		log.Printf("Account doesn't exist")
		return &pb.DepositResponse{
			Transaction: req.GetTransaction(),
		}, err
	}

	amount := account.Balance + req.GetTransaction().GetAmount()
	if amount < 0 {
		log.Printf("the balance is not enough to make this transaction")
		return &pb.DepositResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Account:     toPB(account),
		}, nil
	}

	account.Balance = amount
	if err := db.Updates(account).Error; err != nil {
		log.Println("Failed to withdraw")
		return &pb.DepositResponse{
			Message: "Failed",
			Account: toPB(account),
		}, err
	}

	res := &pb.DepositResponse{
		Message:     "OK",
		Transaction: req.GetTransaction(),
		Account:     toPB(account),
	}
	return res, nil
}

func (s *server) Acc_Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	db := s.connector

	from := req.GetTransaction().GetFrom()
	to := req.GetTransaction().GetTo()
	// fmt.Println(from, to)

	var sender *model.AccountORM
	var receiver *model.AccountORM

	if req.GetTransaction().GetAmount() < 0 {
		log.Printf("Please enter an amount greater than 0")
		return &pb.TransferResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Sender:      toPB(sender),
			Receiver:    toPB(receiver),
		}, nil
	}

	err1 := db.Where(`id =?`, from).First(&sender).Error
	err2 := db.Where(`id =?`, to).First(&receiver).Error
	if err1 != nil {
		log.Printf("Sender doesn't exist")
		return &pb.TransferResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
		}, err1
	}

	if err2 != nil {
		log.Printf("Receiver doesn't exist")
		return &pb.TransferResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
		}, err2
	}

	sender.Balance = sender.Balance - req.GetTransaction().GetAmount()
	receiver.Balance = receiver.Balance + req.GetTransaction().GetAmount()

	fmt.Printf("Sender's Balance : %v, Receiver's Balance %v \n", sender.Balance, receiver.Balance)

	if sender.Balance < 0 {
		log.Printf("the balance is not enough to make this transaction")
		return &pb.TransferResponse{
			Message:  "Failed",
			Sender:   toPB(sender),
			Receiver: toPB(receiver),
		}, nil
	}

	var accounts []*model.AccountORM
	accounts = append(accounts, sender)
	accounts = append(accounts, receiver)

	// fmt.Println(accounts)

	if err := db.Save(accounts).Error; err != nil {
		log.Println("Failed to transfer")
		return &pb.TransferResponse{
			Message:     "Failed",
			Transaction: req.GetTransaction(),
			Sender:      toPB(sender),
			Receiver:    toPB(receiver),
		}, err
	}

	res := &pb.TransferResponse{
		Message:     "OK",
		Transaction: req.GetTransaction(),
		Sender:      toPB(sender),
		Receiver:    toPB(receiver),
	}
	return res, nil
}

func main() {
	fmt.Println("Welcome to the server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to listen the port", err)
	}
	s := grpc.NewServer()

	// GORM
	db := initDB()

	pb.RegisterAccountServiceServer(s, &server{connector: db})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
