package model

import (
	"github.com/rs/xid"
)

type UserID = int64
type AccountID = string

type User struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	TelegramID  int64                 `json:"telegramId"`
	Accounts    map[AccountID]Account `json:"accounts"`
	TradeNotify bool                  `json:"tradeNotify"`
}

func NewUser(name string, telegramID int64, accounts []Account) User {
	return User{
		ID:         xid.New().String(),
		Name:       name,
		TelegramID: telegramID,
		Accounts:   map[AccountID]Account{},
	}
}

type Account struct {
	Name      string  `json:"name"`
	ApiKey    string  `json:"apiKey"`
	SecretKey string  `json:"secretKey"`
	Base      float64 `json:"base"`
}

func NewAccount(name string, apiKey string, secretKey string) Account {
	return Account{
		Name:      name,
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}
}
