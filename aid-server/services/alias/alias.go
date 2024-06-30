package alias

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// DB represents the user database
type DB struct {
	db *sql.DB
}

// NewDB creates a new alias database
func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	udb := &DB{db: db}
	if err := udb.initialize(); err != nil {
		panic(err.Error())
	}

	return udb, nil
}

func (d *DB) initialize() error {
	_, err := d.db.Exec(`
        CREATE TABLE IF NOT EXISTS Users (
            Uid TEXT PRIMARY KEY,
            Name TEXT NOT NULL,
            Pin TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS LoginRecords (
            Id INTEGER PRIMARY KEY AUTOINCREMENT,
            Uid TEXT NOT NULL,
            IP TEXT NOT NULL,
            Browser TEXT NOT NULL,
            LoginTime TEXT NOT NULL,
            FOREIGN KEY (Uid) REFERENCES Users(Uid)
        );`)
	return err
}

// AddUser adds a new user to the database
func (d *DB) AddUser(uid, name, pin string) error {
	_, err := d.db.Exec("INSERT INTO Users (Uid, Name, Pin) VALUES (?, ?, ?)", uid, name, pin)
	return err
}

// ValidateUser checks if the given name and pin are valid
func (d *DB) ValidateUser(name, pin string) ([]string, error) {
	rows, err := d.db.Query("SELECT Uid FROM Users WHERE Name = ? AND Pin = ?", name, pin)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic("close with error")
		}
	}(rows)

	var validUids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		validUids = append(validUids, uid)
	}
	return validUids, nil
}

// AddLoginRecord adds a new login record for a user
func (d *DB) AddLoginRecord(uid, ip, browser string) error {
	_, err := d.db.Exec(
		"INSERT INTO LoginRecords (Uid, IP, Browser, LoginTime) VALUES (?, ?, ?, ?)",
		uid, ip, browser, time.Now().UTC().Format(time.RFC3339))
	return err
}

// LoginRecord represents a single login record
type LoginRecord struct {
	IP        string
	Browser   string
	LoginTime time.Time
}

// GetUserLoginHistory retrieves the most recent login record for a user
func (d *DB) GetUserLoginHistory(uid string) (*LoginRecord, error) {
	row := d.db.QueryRow(
		"SELECT IP, Browser, LoginTime FROM LoginRecords WHERE Uid = ? ORDER BY LoginTime DESC LIMIT 1",
		uid)

	var record LoginRecord
	var loginTimeStr string
	err := row.Scan(&record.IP, &record.Browser, &loginTimeStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No login history found
		}
		return nil, err
	}

	record.LoginTime, err = time.Parse(time.RFC3339, loginTimeStr)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}
