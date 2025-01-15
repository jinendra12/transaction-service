package model

import "time"

type Transaction struct {
	ID        int64        `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Amount    float64      `gorm:"type:decimal(20,2);not null" json:"amount"`
	Type      string       `gorm:"type:varchar(255);not null;index" json:"type"`
	ParentID  *int64       `gorm:"index" json:"parent_id,omitempty"`
	Parent    *Transaction `gorm:"foreignKey:ParentID" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionRequest struct {
	Amount   float64 `json:"amount" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	ParentID *int64  `json:"parent_id,omitempty"`
}

type TransactionResponse struct {
	Amount   float64 `json:"amount"`
	Type     string  `json:"type"`
	ParentID *int64  `json:"parent_id,omitempty"`
}

type SumResponse struct {
	Sum float64 `json:"sum"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type NodeTraversal struct {
	Visited map[int64]bool
}
