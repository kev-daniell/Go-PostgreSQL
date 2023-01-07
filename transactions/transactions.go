package transaction

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type ApiService interface {
	MakeTransaction(w http.ResponseWriter, r *http.Request)
}

type Transaction struct {
	gorm.Model
	Method string
	Amount int
	Item   string
	UserID int
}

type txn struct {
	db *gorm.DB
}

func NewTxn(db *gorm.DB) *txn {
	return &txn{db}
}

func (t *txn) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		txnRow Transaction
	)

	err := json.NewDecoder(r.Body).Decode(&txnRow)
	if err != nil {
		log.Println("error occurred decoding", err)
		return
	}

	createdTxn := t.db.Create(&txnRow)
	if createdTxn.Error != nil {
		log.Println("error occurred creating", createdTxn.Error)
	}

	err = json.NewEncoder(w).Encode(&createdTxn)
	if err != nil {
		return
	}
}

func (t *txn) GetTransaction(w http.ResponseWriter, r *http.Request) {
	var txns []Transaction

	t.db.Find(&txns)

	err := json.NewEncoder(w).Encode(&txns)
	if err != nil {
		return
	}
}
