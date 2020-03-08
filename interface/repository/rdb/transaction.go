package rdb

import (
	"github.com/jinzhu/gorm"
	"github.com/x-color/vue-trello/usecase"
)

// TransactionManager is DB transaction manager.
type TransactionManager struct {
	db *gorm.DB
}

func (m *TransactionManager) BeginTransaction(on bool) usecase.Transaction {
	if on {
		return &Transaction{
			db: m.db.Begin(),
			on: true,
		}
	}

	return &Transaction{
		db: m.db,
		on: false,
	}
}

func newTransactionManager(db *gorm.DB) TransactionManager {
	return TransactionManager{
		db: db,
	}
}

type Transaction struct {
	db *gorm.DB
	on bool
}

func (tx *Transaction) DB() interface{} {
	return interface{}(tx.db)
}

func (tx *Transaction) Commit() {
	if tx.on {
		tx.db.Commit()
	}
}

func (tx *Transaction) Rollback() {
	if tx.on {
		tx.db.Rollback()
	}
}
