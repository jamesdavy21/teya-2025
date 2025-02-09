package client

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jamesdavy21/teya-2025/internal/application"
)

type TransactionManager interface {
	AddDeposit(accountID uuid.UUID, amount float64) (*application.Transaction, error)
	AddWithdrawal(accountID uuid.UUID, amount float64) (*application.Transaction, error)
	GetTransactions(accountID uuid.UUID, page, limit int) ([]application.Transaction, int, error)
	GetAccount(accountID uuid.UUID) (*application.Account, error)
}

type TransactionClient struct {
	transactionManager TransactionManager
}

func NewTransactionClient(transactionManager TransactionManager) *TransactionClient {
	return &TransactionClient{
		transactionManager: transactionManager,
	}
}

type Deposit struct {
	Amount float64 `json:"amount"`
}

func (t *TransactionClient) HandleDeposit(c *gin.Context) {
	accountIDString := c.Param("id")

	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		slog.Error("Error parsing account ID", slog.String("accountID", accountIDString))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var deposit Deposit
	err = c.BindJSON(&deposit)
	if err != nil {
		slog.Error("Error parsing deposit", slog.String("accountID", accountIDString), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if deposit.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	transaction, err := t.transactionManager.AddDeposit(accountID, deposit.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

type Withdrawal struct {
	Amount float64 `json:"amount"`
}

func (t *TransactionClient) HandleWithdrawal(c *gin.Context) {
	accountIDString := c.Param("id")

	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		slog.Error("Error parsing account ID", slog.String("accountID", accountIDString))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody Withdrawal
	err = c.BindJSON(&reqBody)
	if err != nil {
		slog.Error("Error parsing withdrawal", slog.String("accountID", accountIDString), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqBody.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	transaction, err := t.transactionManager.AddWithdrawal(accountID, reqBody.Amount)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "create account by making a valid deposit first"})
		case errors.Is(err, application.ErrNotEnoughFunds):
			c.JSON(http.StatusBadRequest, gin.H{"error": "not enough funds in account"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

type GetTransactionsRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (t *TransactionClient) HandleGetTransactions(c *gin.Context) {
	accountIDString := c.Param("id")

	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		slog.Error("Error parsing account ID", slog.String("accountID", accountIDString))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody GetTransactionsRequest
	err = c.ShouldBindBodyWithJSON(&reqBody)
	if err != nil && !errors.Is(err, io.EOF) {
		slog.Error("Error parsing withdrawal", slog.String("accountID", accountIDString), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqBody.Page <= 0 {
		reqBody.Page = 0
	}

	if reqBody.Limit <= 0 || reqBody.Limit > 25 {
		reqBody.Limit = 25
	}

	transactions, nextPage, err := t.transactionManager.GetTransactions(accountID, reqBody.Page, reqBody.Limit)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrAccountNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "create account by making a valid deposit first"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions, "next_page": nextPage})
}

func (t *TransactionClient) HandleGetAccount(c *gin.Context) {
	accountIDString := c.Param("id")

	accountID, err := uuid.Parse(accountIDString)
	if err != nil {
		slog.Error("Error parsing account ID", slog.String("accountID", accountIDString))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := t.transactionManager.GetAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}
