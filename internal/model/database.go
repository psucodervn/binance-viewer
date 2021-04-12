package model

import (
	"sync"
)

type Database struct {
	Users map[UserID]User `json:"users"`
	mu    sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		Users: map[UserID]User{},
	}
}

func (d *Database) Renew() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Users = nil
}

func (d *Database) AddUser(user User) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Users[user.TelegramID]; ok {
		return nil
	}
	d.Users[user.TelegramID] = user
	return nil
}

func (d *Database) FindUser(id int64) (User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for i := range d.Users {
		if d.Users[i].TelegramID == id {
			u := d.Users[i]
			return u, nil
		}
	}
	return User{}, ErrNotFound
}

func (d *Database) AddAccount(u User, acc Account) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.Users[u.TelegramID]; !ok {
		return ErrNotFound
	}
	if _, ok := d.Users[u.TelegramID].Accounts[acc.ApiKey]; ok {
		return nil
	}
	d.Users[u.TelegramID].Accounts[acc.ApiKey] = acc
	return nil
}
