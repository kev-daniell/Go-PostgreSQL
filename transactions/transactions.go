package transaction

import (
	"Postgres/types"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type TxnApi interface {
	MakeTransaction(w http.ResponseWriter, r *http.Request)
	GetAllTransactions(w http.ResponseWriter, r *http.Request)
	GetTransaction(w http.ResponseWriter, r *http.Request)
	UpdateTransaction(w http.ResponseWriter, r *http.Request)
	DeleteTransaction(w http.ResponseWriter, r *http.Request)
}

type txn struct {
	db *gorm.DB
}

func NewTxn(db *gorm.DB) *txn {
	return &txn{db}
}

func (t *txn) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var txnRow types.Transaction

	err := json.NewDecoder(r.Body).Decode(&txnRow)
	if err != nil {
		log.Println("error occurred decoding", err)
		return
	}

	resp := t.db.Create(&txnRow)
	if resp.Error != nil {
		log.Println("error occurred creating", resp.Error)
	}

	err = json.NewEncoder(w).Encode(&txnRow)
	if err != nil {
		return
	}
}

func (t *txn) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []types.Transaction

	t.db.Find(&transactions)

	err := json.NewEncoder(w).Encode(&transactions)
	if err != nil {
		return
	}
}

func (t *txn) GetTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction types.Transaction

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("id is missing in GetUser")
		log.Println(err)
		w.WriteHeader(404)
	}

	resp := t.db.First(&transaction, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
	}

	err = json.NewEncoder(w).Encode(&transaction)
	if err != nil {
		return
	}
}

func (t *txn) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		newTxn types.Transaction
		oldTxn types.Transaction
	)

	err := json.NewDecoder(r.Body).Decode(&newTxn)
	if err != nil {
		log.Println("error occurred decoding", err)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("id is missing in GetUser")
		log.Println(err)
		w.WriteHeader(404)
	}

	resp := t.db.First(&oldTxn, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
	}

	resp = t.db.Model(&oldTxn).Updates(newTxn)
	if resp.Error != nil {
		log.Println("error occurred creating", resp.Error)
	}

	err = json.NewEncoder(w).Encode(&newTxn)
	if err != nil {
		return
	}
}

func (t *txn) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	var txnToDelete types.Transaction

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("id is missing in GetUser")
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	resp := t.db.First(&txnToDelete, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
		return
	}

	resp = t.db.Delete(&txnToDelete)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
		return
	}

	err = json.NewEncoder(w).Encode(&txnToDelete)
	if err != nil {
		return
	}
}
