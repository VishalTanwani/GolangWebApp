package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn" // this is making postgres connection
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the data base connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpneDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

//ConnectSQL creates database pool connection for postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpneDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

//testDB tries to ping database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

//NewDatabase creates a new database connection for the application
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
