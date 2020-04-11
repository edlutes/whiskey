package main

import (
	"database/sql"
)

type Store interface {
	enterWhiskey(whiskey *Whiskey) error
	getWhiskeys() ([]*Whiskey, error)
}
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) enterWhiskey(whiskey *Whiskey) error {

	_, err := store.db.Query("INSERT INTO whiskey_collection(name, description) VALUES ($1,$2)", whiskey.Name, whiskey.Description)
	return err
}

func (store *dbStore) getWhiskeys() ([]*Whiskey, error) {

	rows, err := store.db.Query("SELECT name, description from whiskey_collection")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	whiskeys := []*Whiskey{}
	for rows.Next() {
		whiskey := &Whiskey{}
		if err := rows.Scan(&whiskey.Name, &whiskey.Description); err != nil {
			return nil, err
		}

		whiskeys = append(whiskeys, whiskey)
	}
	return whiskeys, nil
}

var store Store

func InitStore(s Store) {
	store = s
}
