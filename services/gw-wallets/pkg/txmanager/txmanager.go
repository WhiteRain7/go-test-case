package txmanager

import (
	"github.com/jmoiron/sqlx"
)

type TXManager struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewTXManager(db *sqlx.DB) *TXManager {
	return &TXManager{nil, db}
}

func (txm *TXManager) Get() sqlx.ExtContext {
	if txm.tx == nil {
		return txm.db
	}
	return txm.tx
}

func (txm *TXManager) BeginTx() error {
	tx, err := txm.db.Beginx()
	if err != nil {
		return err
	}
	txm.tx = tx
	return nil
}

func (txm *TXManager) Commit() error {
	defer func() {
		txm.tx = nil
	}()
	if txm.tx == nil {
		return nil
	}
	return txm.tx.Commit()
}

func (txm *TXManager) Rollback() error {
	defer func() {
		txm.tx = nil
	}()
	if txm.tx == nil {
		return nil
	}
	return txm.tx.Rollback()
}
