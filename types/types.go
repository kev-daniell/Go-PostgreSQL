package types

import "gorm.io/gorm"

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
	Transactions []Transaction
}
