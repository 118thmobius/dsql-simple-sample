package main

import (
	"context"
	"dsql-simple-sample/domain"
	"dsql-simple-sample/infrastructure"
	"dsql-simple-sample/infrastructure/db"
	"dsql-simple-sample/service"
	"dsql-simple-sample/usecase"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var accountUseCase usecase.AccountUseCase

var rootCmd = &cobra.Command{
	Use: "dsql-simple-sample",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var transferCmd = &cobra.Command{
	Use:  "transfer",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fromId := args[0]
		toId := args[1]
		amount := args[2]

		amountInt, err := strconv.Atoi(amount)
		if err != nil {
			fmt.Println("amount must be integer.")
		}

		err = accountUseCase.Transfer(context.Background(), fromId, toId, amountInt)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Transfer", amountInt, "from", fromId, "to", toId, "amount", amountInt)
	},
}

var showCmd = &cobra.Command{
	Use:  "show",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		account, err := accountUseCase.GetAccountByID(context.Background(), id)
		if err != nil {
			fmt.Println(err)
		} else {
			printAccount(account)
		}
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
	rootCmd.AddCommand(showCmd)
}

func main() {
	ctx := context.Background()
	region := os.Getenv("AWS_REGION")
	clusterEndpoint := os.Getenv("AWS_CLUSTER_ENDPOINT")
	if region == "" || clusterEndpoint == "" {
		panic("AWS_REGION or AWS_CLUSTER_ENDPOINT is not set.")
	}

	pool, err := db.GetPool(ctx, region, clusterEndpoint)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	txManager := db.NewTransactionManager(pool)
	accountRepository := infrastructure.NewAccountRepository()
	transactionRepository := infrastructure.NewTransactionRepository()
	txDomainService := service.NewTransactionDomainService()
	accountUseCase = usecase.NewAccountUseCase(txManager, accountRepository, transactionRepository, txDomainService)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func printAccount(account *domain.Account) {
	fmt.Println("UserID:", account.UserId, "; City:", account.City, "; Balance:", account.Balance)
}
