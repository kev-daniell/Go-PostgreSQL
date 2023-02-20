package transaction

import (
	"Postgres/types"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type TxnApi interface {
	MakeTransaction(w http.ResponseWriter, r *http.Request)
	GetTransaction(w http.ResponseWriter, r *http.Request)
}

type txn struct {
	db *gorm.DB
}

func NewTxn(db *gorm.DB) *txn {
	return &txn{db}
}

func (t *txn) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		txnRow types.Transaction
		usr    types.User
	)

	err := json.NewDecoder(r.Body).Decode(&txnRow)
	if err != nil {
		log.Println("error occurred decoding", err)
		return
	}

	resp := t.db.Create(&txnRow)
	if resp.Error != nil {
		log.Println("error occurred creating", resp.Error)
	}

	t.db.Where("id = ?", txnRow.UserID).Find(&usr)
	usr.Transactions = append(usr.Transactions, txnRow)
	t.db.Save(&usr)

	err = json.NewEncoder(w).Encode(&txnRow)
	if err != nil {
		return
	}
}

func (t *txn) GetTransaction(w http.ResponseWriter, r *http.Request) {
	var transactions []types.Transaction

	t.db.Find(&transactions)

	err := json.NewEncoder(w).Encode(&transactions)
	if err != nil {
		return
	}
}
