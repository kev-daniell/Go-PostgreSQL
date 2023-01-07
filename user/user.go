package user

import (
	transaction "Postgres/transactions"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type ApiService interface {
	MakeUser(w http.ResponseWriter, r *http.Request)
}

type Transaction struct {
	gorm.Model
	Method string
	Amount int
	Item   string
	UserID int
}

type User struct {
	gorm.Model
	Name         string
	Transactions []transaction.Transaction
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *user {
	return &user{db}
}

func (u *user) MakeUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error decoding", err)
		return
	}

	createdUser := u.db.Create(&user)
	if createdUser.Error != nil {
		log.Println("error making ", user, createdUser.Error)
		return
	}

	err = json.NewEncoder(w).Encode(&createdUser)
	if err != nil {
		log.Println("error")
		return
	}

}

func (u *user) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	var txns []Transaction

	u.db.Find(&users)
	u.db.Model(&users).Related(&txns)
	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		return
	}
}
