package infra_database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Database interface {
	Select(query string, params ...any) (*sql.Rows, error)
	SelectOne(query string, params ...any) *sql.Row
	Exec(query string, params ...any) error
	Transaction(operation Operation) error
	QueryBuilder(table string) *QueryBuilder
	Connect()
}

type Operation func(tx *sql.Tx) error

type MySQLDatabase struct {
	config mysql.Config
	db     *sql.DB
}

func (database *MySQLDatabase) Connect() {
	connetion, err := sql.Open("mysql", database.config.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	database.db = connetion

	pingErr := database.db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func (database *MySQLDatabase) Select(query string, params ...any) (*sql.Rows, error) {
	rows, err := database.db.Query(query, params...)
	return rows, err
}

func (database *MySQLDatabase) Transaction(operation Operation) error {
	tx, err := database.db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	transaction_err := operation(tx)

	if transaction_err != nil {
		tx.Rollback()
		return errors.New("could not complete transaction")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("could not commit transaction")
	}

	return nil
}

func (database *MySQLDatabase) SelectOne(query string, params ...any) *sql.Row {
	row := database.db.QueryRow(query, params...)
	return row
}

func (database *MySQLDatabase) QueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

func (database *MySQLDatabase) Exec(query string, params ...any) error {
	_, err := database.db.Exec(query, params...)
	return err
}

func NewMySQLDatabase(config MySQLDatabaseConfig) *MySQLDatabase {
	dbConfig := mysql.Config{
		User:      config.User,
		Passwd:    config.Password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%d", config.Host, config.Port),
		DBName:    config.Name,
		ParseTime: true,
	}

	return &MySQLDatabase{config: dbConfig, db: nil}
}

type MySQLDatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
}
