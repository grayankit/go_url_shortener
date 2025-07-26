package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DB struct {
	conn *sql.DB
}
type URLRecord struct {
	Code      string
	LongURL   string
	CreatedAt string
}

func New(path string) *DB {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open database %v", err)
	}

	_, err = conn.Exec(`
	CREATE TABLE IF NOT EXISTS urls (
		code TEXT PRIMARY KEY,
		long_url TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return &DB{conn: conn}
}
func (db *DB) Save(code, longURL string) error {
	_, err := db.conn.Exec(`INSERT INTO urls (code, long_url) VALUES (?, ?)`, code, longURL)
	return err
}
func (db *DB) GetLongURL(code string) (string, error) {
	row := db.conn.QueryRow(`SELECT long_url FROM urls WHERE code = ?`, code)
	var url string
	err := row.Scan(&url)
	return url, err
}
func (db *DB) GetCodeByURL(longURL string) (string, error) {
	row := db.conn.QueryRow(`SELECT code FROM urls WHERE long_url = ?`, longURL)
	var code string
	err := row.Scan(&code)
	return code, err
}
func (db *DB) GetAll() ([]URLRecord, error) {
	rows, err := db.conn.Query(`SELECT code, long_url, created_at FROM urls ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URLRecord
	for rows.Next() {
		var u URLRecord
		err := rows.Scan(&u.Code, &u.LongURL, &u.CreatedAt)
		if err != nil {
			continue
		}
		urls = append(urls, u)
	}
	return urls, nil
}
func (db *DB) Delete(code string) error {
	_, err := db.conn.Exec(`DELETE FROM urls WHERE code = ?`, code)
	return err
}
