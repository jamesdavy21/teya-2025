package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jamesdavy21/teya-2025/client"
	"github.com/jamesdavy21/teya-2025/internal/application"
	"github.com/jamesdavy21/teya-2025/internal/store/memory"
)

func main() {

	router := gin.Default()

	store := memory.NewInMemoryStore()
	app := application.NewTransactionManager(store)

	transactionClient := client.NewTransactionClient(app)

	router.POST("/account/:id/deposit", transactionClient.HandleDeposit)
	router.POST("/account/:id/withdrawal", transactionClient.HandleWithdrawal)
	router.GET("/account/:id/transactions", transactionClient.HandleGetTransactions)
	router.GET("/account/:id", transactionClient.HandleGetAccount)

	err := router.Run(":8080")
	if err != nil {
		slog.Error("Error running server", slog.String("error", err.Error()))
	}
}
