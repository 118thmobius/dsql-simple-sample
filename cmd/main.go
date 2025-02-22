package main

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/service"
	"dsql-simple-sample/usecase"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	region := os.Getenv("AWS_REGION")
	clusterEndpoint := os.Getenv("AWS_CLUSTER_ENDPOINT")
	pool, err := db.GetPool(ctx, region, clusterEndpoint)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	txManager := db.NewTransactionManager(pool)
	accountRepository := infrastructure.NewAccountRepository()
	transactionRepository := infrastructure.NewTransactionRepository()
	txDomainService := service.NewTransactionDomainService()
	accountUseCase := usecase.NewAccountUseCase(txManager, accountRepository, transactionRepository, txDomainService)

	account, err := accountUseCase.GetAccountByID(ctx, "Alice")
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] Before transfer:")
	printAccount(account)
	fmt.Println("")

	fmt.Println("[*] Begin transfer...")
	fmt.Println("Transfer 500 from Alice to Bob.")
	if err := accountUseCase.Transfer(ctx, "Alice", "Bob", 500); err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("Transfer completed.")
	}

	account, err = accountUseCase.GetAccountByID(ctx, "Alice")
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] After transfer:")
	printAccount(account)
}

func printAccount(account *domain.Account) {
	fmt.Println("UserID:", account.UserId, "; City:", account.City, "; Balance:", account.Balance)
}
