package model

import (
	"encoding/json"
	"sync"
)

type Database struct {
	users map[UserID]User
	mu    sync.RWMutex
}

type tmpJson struct {
	Users map[UserID]User `json:"users"`
}

func (d *Database) UnmarshalJSON(bytes []byte) error {
	var tmp tmpJson
	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return err
	}
	d.mu.Lock()
	d.users = tmp.Users
	d.mu.Unlock()
	return nil
}

func (d *Database) MarshalJSON() ([]byte, error) {
	tmp := tmpJson{Users: d.Users()}
	return json.Marshal(tmp)
}

func (d *Database) Users() map[UserID]User {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.users
}

func NewDatabase() *Database {
	return &Database{
		users: map[UserID]User{},
	}
}

func (d *Database) Renew() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.users = nil
}

func (d *Database) AddUser(user User) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.users[user.TelegramID]; ok {
		return nil
	}
	d.users[user.TelegramID] = user
	return nil
}

func (d *Database) FindUser(id int64) (User, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for i := range d.users {
		if d.users[i].TelegramID == id {
			u := d.users[i]
			return u, nil
		}
	}
	return User{}, ErrNotFound
}

func (d *Database) AddAccount(u User, acc Account) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.users[u.TelegramID]; !ok {
		return ErrNotFound
	}
	if _, ok := d.users[u.TelegramID].Accounts[acc.ApiKey]; ok {
		return nil
	}
	d.users[u.TelegramID].Accounts[acc.ApiKey] = acc
	return nil
}
