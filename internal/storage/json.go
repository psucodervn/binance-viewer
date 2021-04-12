package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	DefaultFilePath = "data/data.json"
)

func LoadOrCreate(db interface{}) error {
	b, err := ioutil.ReadFile(DefaultFilePath)
	if os.IsNotExist(err) {
		err = Save(db)
		return err
	}
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, db); err != nil {
		return err
	}
	return nil
}

func Save(db interface{}) error {
	b, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(DefaultFilePath, b, 0600)
}
