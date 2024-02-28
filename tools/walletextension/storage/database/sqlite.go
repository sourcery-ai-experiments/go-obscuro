package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	obscurocommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	_ "github.com/mattn/go-sqlite3" // sqlite driver for sql.Open()
	common "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type SqliteDatabase struct {
	db *sql.DB
}

func NewSqliteDatabase(dbPath string) (*SqliteDatabase, error) {
	// load the db file
	dbFilePath, err := createOrLoad(dbPath)
	if err != nil {
		return nil, err
	}

	// open the db
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		fmt.Println("Error opening database: ", err)
		return nil, err
	}

	// enable foreign keys in sqlite
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// create users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		user_id binary(20) PRIMARY KEY,
		private_key binary(32)
	);`)

	if err != nil {
		return nil, err
	}

	// create accounts table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		user_id binary(20),
		account_address binary(20),
		signature binary(65),
		signature_type int,
    	FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE
	);`)

	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{db: db}, nil
}

func (s *SqliteDatabase) AddUser(userID []byte, privateKey []byte) error {
	stmt, err := s.db.Prepare("INSERT OR REPLACE INTO users(user_id, private_key) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDatabase) DeleteUser(userID []byte) error {
	stmt, err := s.db.Prepare("DELETE FROM users WHERE user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDatabase) GetUserPrivateKey(userID []byte) ([]byte, error) {
	var privateKey []byte
	err := s.db.QueryRow("SELECT private_key FROM users WHERE user_id = ?", userID).Scan(&privateKey)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given userID
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}

	return privateKey, nil
}

func (s *SqliteDatabase) AddAccount(userID []byte, accountAddress []byte, signature []byte, signatureType int) error {
	stmt, err := s.db.Prepare("INSERT INTO accounts(user_id, account_address, signature, signature_type) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, accountAddress, signature, signatureType)
	if err != nil {
		return err
	}

	return nil
}

func (s *SqliteDatabase) GetAccounts(userID []byte) ([]common.AccountDB, error) {
	rows, err := s.db.Query("SELECT account_address, signature, signature_type FROM accounts WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []common.AccountDB
	for rows.Next() {
		var account common.AccountDB
		if err := rows.Scan(&account.AccountAddress, &account.Signature, &account.SignatureType); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *SqliteDatabase) GetAllUsers() ([]common.UserDB, error) {
	rows, err := s.db.Query("SELECT user_id, private_key FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []common.UserDB
	for rows.Next() {
		var user common.UserDB
		err = rows.Scan(&user.UserID, &user.PrivateKey)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func createOrLoad(dbPath string) (string, error) {
	// If path is empty we create a random throwaway temp file, otherwise we use the path to the database
	if dbPath == "" {
		tempDir := filepath.Join("/tmp", "obscuro_gateway", obscurocommon.RandomStr(8))
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory: ", tempDir, err)
			return "", err
		}
		dbPath = filepath.Join(tempDir, "gateway_databse.db")
	} else {
		dir := filepath.Dir(dbPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			fmt.Println("Error creating directories:", err)
			return "", err
		}
	}

	return dbPath, nil
}
