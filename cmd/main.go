package main

import (
	"context"
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

	accountRepository := infrastructure.NewAccountRepository()
	transactionRepository := infrastructure.NewTransactionRepository()
	txDomainService := service.NewTransactionDomainService()
	accountUseCase := usecase.NewAccountUseCase(pool, accountRepository, transactionRepository, txDomainService)

	if err := accountUseCase.Transfer(ctx, "Alice", "Bob", 500); err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		println("success")
	}

}
