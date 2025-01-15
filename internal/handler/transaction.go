package handler

import (
	"net/http"
	"strconv"
	"transaction-service/internal/model"
	"transaction-service/internal/service"
	"transaction-service/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) RegisterRoutes(router *gin.Engine) {
	transactionGroup := router.Group("/transactionservice")
	{
		transactionGroup.PUT("/transaction/:id", h.CreateTransaction)
		transactionGroup.GET("/transaction/:id", h.GetTransaction)
		transactionGroup.GET("/types/:type", h.GetTransactionsByType)
		transactionGroup.GET("/sum/:id", h.GetTransactionSum)
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	var req model.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = h.service.CreateTransaction(id, &req)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.IsInvalidData(err):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, model.StatusResponse{
		Status: "ok",
	})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	transaction, err := h.service.GetTransaction(id)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetTransactionsByType(c *gin.Context) {
	txType := c.Param("type")

	ids, err := h.service.GetTransactionsByType(txType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, ids)
}

func (h *TransactionHandler) GetTransactionSum(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid transaction ID",
		})
		return
	}

	sum, err := h.service.CalculateTransactionSum(id)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, model.SumResponse{
		Sum: sum,
	})
}
