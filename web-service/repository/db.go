package repository

import (
	"database/sql"
	"fmt"
)

// type IRepository interface {
// 	Create() error
// 	ReadAll() ([]interface{}, error)
// 	Read(string) (interface{}, error)
// 	Update(string) error
// 	Delete(string) error
// }

func GetConnection(host string, port int, user, password, dbName string) (*sql.DB, error) {
	if password == "" {
		password = `''`
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	// connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbName)
	// fmt.Println("COnnection: " + connStr + "\n")
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}
	return database, nil
}
