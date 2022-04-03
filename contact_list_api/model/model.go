package main

import (
	"database/sql"
)

type contact struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func (c *contact) getContact(db *sql.DB) error {
	return db.QueryRow("SELECT name, number FROM contact WHERE id=$1",
		c.ID).Scan(&c.Name, &c.Number)
}

func (c *contact) updateContact(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE contact SET name=$1, number=$2 WHERE id=$3",
			c.Name, c.Number, c.ID)

	return err
}

func (c *contact) deleteContact(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM contact WHERE id=$1", c.ID)

	return err
}

func (c *contact) createContact(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO contact(name, number) VALUES($1, $2) RETURNING id",
		c.Name, c.Number).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getAllContacts(db *sql.DB, start, count int) ([]contact, error) {
	rows, err := db.Query(
		"SELECT id, name,  number FROM contacts LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contacts := []contact{}

	for rows.Next() {
		var c contact
		if err := rows.Scan(&c.ID, &c.Name, &c.Number); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}

	return contacts, nil
}
