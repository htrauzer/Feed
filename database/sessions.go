package database

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

func CreateSessionsTable(db *sql.DB) error {
	sessionsTable := `CREATE TABLE Sessions (
		id TEXT PRIMARY KEY,		
		userId TEXT,
		created_at TEXT
	  );`

	log.Println("Creating Sessions table...")
	statement, err := db.Prepare(sessionsTable)
	if err != nil {
		return err
	}
	statement.Exec()
	log.Println("Table created")
	return nil
}

func NewSession(db *sql.DB, cookieValue, usr string) error {
	q := `INSERT INTO Sessions (id, userId, created_at) VALUES (?,?,?);`
	statement, err := db.Prepare(q)
	if err != nil {
		return err
	}
	timenow := time.Now().String()
	_, err = statement.Exec(cookieValue, usr, timenow)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSession(db *sql.DB, cookie *http.Cookie) error { // by uuid
	printQuery := `SELECT created_at FROM Sessions WHERE id = ?`
	row, _ := db.Query(printQuery, cookie.Value)
	defer row.Close()
	for row.Next() {
		var created string
		row.Scan(&created)
	}

	q := `DELETE FROM Sessions WHERE id = ?;`
	statement, err := db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = statement.Exec(cookie.Value)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSessionByUserName(db *sql.DB, loginUserName string) error { // when logging in from other browser
	printQuery := `SELECT created_at FROM Sessions WHERE userId = ?`
	row, _ := db.Query(printQuery, loginUserName)
	defer row.Close()
	for row.Next() {
		var created string
		row.Scan(&created)
	}

	q := `DELETE FROM Sessions WHERE userId = ?;`
	statement, err := db.Prepare(q)
	if err != nil {
		return err
	}
	_, err = statement.Exec(loginUserName)
	if err != nil {
		return err
	}
	return nil
}

func HasSession(db *sql.DB, cookieValue string) (bool, string, error) {
	row, err := db.Query("SELECT id, userId FROM Sessions")
	if err != nil {
		return false, "", err
	}
	defer row.Close()

	for row.Next() {
		var id string
		var userId string
		row.Scan(&id, &userId)
		if id == cookieValue {
			return true, userId, nil
		}
	}
	return false, "", nil
}
