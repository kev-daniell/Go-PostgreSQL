package user

import (
	"Postgres/types"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type ApiService interface {
	MakeUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *user {
	return &user{db}
}

func (u *user) MakeUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("error decoding", err)
		return
	}

	resp := u.db.Create(&user)
	if resp.Error != nil {
		log.Println("error making ", user, resp.Error)
		return
	}

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		log.Println("error")
		return
	}

}

func (u *user) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []types.User

	u.db.Find(&users)
	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		return
	}
}

func (u *user) GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		usr          types.User
		transactions []types.Transaction
	)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("id is missing in GetUser")
		log.Println(err)
		w.WriteHeader(404)
	}

	resp := u.db.First(&usr, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
	}

	resp = u.db.Where("user_id = ?", id).Find(&transactions)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
	}

	usr.Transactions = transactions

	err = json.NewEncoder(w).Encode(&usr)
	if err != nil {
		return
	}
}

func (u *user) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		newUser types.User
		oldUser types.User
	)

	err := json.NewDecoder(r.Body).Decode(&newUser)
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

	resp := u.db.First(&oldUser, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
	}

	resp = u.db.Model(&oldUser).Updates(newUser)
	if resp.Error != nil {
		log.Println("error occurred creating", resp.Error)
	}

	err = json.NewEncoder(w).Encode(&newUser)
	if err != nil {
		return
	}

}

func (u *user) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userToDelete types.User

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("id is missing in GetUser")
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	resp := u.db.First(&userToDelete, id)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
		return
	}

	resp = u.db.Delete(&userToDelete)
	if resp.Error != nil {
		log.Println(resp.Error)
		w.WriteHeader(500)
		return
	}

	err = json.NewEncoder(w).Encode(&userToDelete)
	if err != nil {
		return
	}
}
