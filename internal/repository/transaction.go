package repository

import (
	"database/sql"
	"fmt"
	"transaction-service/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Save(transaction *model.Transaction) error
	FindByID(id int64) (*model.Transaction, error)
	FindByType(txType string) ([]int64, error)
	FindAllLinkedTransactions(id int64) ([]*model.Transaction, error)
	WouldCreateCycle(transactionID, newParentID int64) (bool, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Save(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) FindByID(id int64) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByType(txType string) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&model.Transaction{}).Where("type = ?", txType).Pluck("id", &ids).Error
	return ids, err
}

func (r *transactionRepository) FindAllLinkedTransactions(id int64) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	fmt.Println("####### been here , checkign code flow #########")
	query := `
        WITH RECURSIVE transaction_tree AS (
            SELECT id, amount, type, parent_id
            FROM transactions
            WHERE id = ?
            
            UNION
            
            SELECT t.id, t.amount, t.type, t.parent_id
            FROM transactions t
            INNER JOIN transaction_tree tt ON t.parent_id = tt.id
        )
        SELECT * FROM transaction_tree WHERE id != ?;
    `

	err := r.db.Raw(query, id, id).Scan(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) WouldCreateCycle(transactionID, newParentID int64) (bool, error) {
	traversal := &model.NodeTraversal{
		Visited: make(map[int64]bool),
	}

	return r.checkCycleFromNode(transactionID, newParentID, traversal)
}

func (r *transactionRepository) checkCycleFromNode(transactionID, currentID int64, traversal *model.NodeTraversal) (bool, error) {
	if currentID == transactionID {
		return true, nil
	}

	if traversal.Visited[currentID] {
		return false, nil
	}

	traversal.Visited[currentID] = true

	var parentID sql.NullInt64
	err := r.db.Table("transactions").
		Select("parent_id").
		Where("id = ?", currentID).
		Row().
		Scan(&parentID)

	if err != nil {
		return false, err
	}

	if !parentID.Valid {
		return false, nil
	}

	return r.checkCycleFromNode(transactionID, parentID.Int64, traversal)
}
