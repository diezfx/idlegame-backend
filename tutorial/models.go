// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package tutorial

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type NullTransactionType struct {
	TransactionType TransactionType
	Valid           bool // Valid is true if TransactionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransactionType) Scan(value interface{}) error {
	if value == nil {
		ns.TransactionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransactionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransactionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TransactionType), nil
}

type Member struct {
	ID string
}

type Project struct {
	ID   uuid.UUID
	Name string
}

type ProjectMembership struct {
	ProjectID uuid.UUID
	UserID    string
}

type Transaction struct {
	ID              uuid.UUID
	Name            string
	Amount          int32
	SourceID        string
	TransactionType TransactionType
}

type TransactionTarget struct {
	TransactionID uuid.UUID
	UserID        string
}