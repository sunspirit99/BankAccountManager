# BankAccountManager

Steps to run :
1.  Personalize your mysql profile in "config/config.go"
2.  run "go run server/server.go"
3.  then run "go run client/client.go"


Benchmark :
1.  run "go run server/server.go"
2.  run "go test -run [FunctionTest] -v -timeout 1000s"

Results in local :
1.  Create 1000 accounts (using Bidirectional Streaming gRPC) : 96.831s
2.  Deposit 1000 times into an account (using Simple gRPC)  : 108.297s
3.  Perform 1000 transactions transfer One-One (using Simple gRPC) :  93.677s
4.  Transfer 1000 transactions One-Many (using Simple gRPC) : 103.330s
5.  Transfer 1000 transactions Many-One (using Simple gRPC) : 105.942s
