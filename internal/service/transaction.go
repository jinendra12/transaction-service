package service

import (
	"fmt"
	"transaction-service/internal/model"
	"transaction-service/internal/repository"
	"transaction-service/pkg/errors"
)

type TransactionService interface {
	CreateTransaction(id int64, req *model.TransactionRequest) error
	GetTransaction(id int64) (*model.TransactionResponse, error)
	GetTransactionsByType(txType string) ([]int64, error)
	CalculateTransactionSum(id int64) (float64, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) CreateTransaction(id int64, req *model.TransactionRequest) error {

	if req.Amount <= 0 {
		return errors.WrapError(errors.ErrInvalidTransaction, "amount can't be negative")
	}

	if req.ParentID != nil && *req.ParentID == id {
		return errors.WrapError(errors.ErrInvalidTransaction, "self-reference not allowed")
	}

	if req.ParentID != nil {
		_, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			return errors.ErrParentTransactionNotFound
		}

		wouldCycle, err := s.repo.WouldCreateCycle(id, *req.ParentID)
		if err != nil {
			return errors.WrapError(errors.ErrDatabaseOperation, "failed to check for cycles")
		}
		if wouldCycle {
			return errors.WrapError(errors.ErrInvalidTransaction, "operation would create a cycle")
		}
	}

	transaction := &model.Transaction{
		ID:       id,
		Amount:   req.Amount,
		Type:     req.Type,
		ParentID: req.ParentID,
	}

	if err := s.repo.Save(transaction); err != nil {
		return errors.WrapError(errors.ErrDatabaseOperation, fmt.Sprintf("failed to save transaction: %v", err))
	}

	return nil
}

func (s *transactionService) GetTransaction(id int64) (*model.TransactionResponse, error) {
	transaction, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrTransactionNotFound
	}

	return &model.TransactionResponse{
		Amount:   transaction.Amount,
		Type:     transaction.Type,
		ParentID: transaction.ParentID,
	}, nil
}

func (s *transactionService) GetTransactionsByType(txType string) ([]int64, error) {
	if txType == "" {
		return nil, errors.ErrInvalidTransaction
	}

	ids, err := s.repo.FindByType(txType)
	if err != nil {
		return nil, errors.WrapError(errors.ErrDatabaseOperation, "failed to fetch transactions by type")
	}

	return ids, nil
}

func (s *transactionService) CalculateTransactionSum(id int64) (float64, error) {

	rootTx, err := s.repo.FindByID(id)
	if err != nil {
		return 0, errors.ErrTransactionNotFound
	}

	children, err := s.repo.FindAllLinkedTransactions(id)
	if err != nil {
		return 0, errors.WrapError(errors.ErrDatabaseOperation, "failed to get linked transactions")
	}

	sum := rootTx.Amount
	for _, child := range children {
		sum += child.Amount
	}

	return sum, nil
}
