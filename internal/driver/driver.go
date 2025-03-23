package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// ConnectSQL creates database pool for PostgreSQL
func ConnectSql(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err) // panic -> if the database cannot connect, we can't do anything here
	}

	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil

}

// tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
