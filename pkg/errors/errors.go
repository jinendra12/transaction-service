package errors

import (
	"errors"
	"fmt"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")

	ErrInvalidTransaction = errors.New("invalid transaction data")

	ErrParentTransactionNotFound = errors.New("parent transaction not found")

	ErrDatabaseOperation = errors.New("database operation failed")
)

func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrTransactionNotFound)
}

func IsInvalidData(err error) bool {
	return errors.Is(err, ErrInvalidTransaction)
}
