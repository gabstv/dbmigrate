package mh

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
)

var drivers map[string]driver.Driver
var driversMu sync.RWMutex

func Run(fn func(tx Mtx)) {
	//fmt.Println("RUN :: DN", os.Getenv("DBMSESS__DRIVER_NAME"))
	//fmt.Println("RUN :: DSN", os.Getenv("DBMSESS__DATA_SOURCE_NAME"))
	//wdir, _ := os.Getwd()
	//fmt.Println("RUN :: WD", wdir)
	xdb, err := EnvConnect()
	if err != nil {
		stderr(err, 1)
		os.Exit(1)
	}
	defer xdb.Close()
	isoLevel := sql.LevelSerializable
	//TODO: check unsupported in postgres when implementing
	xtx, err := xdb.BeginTxx(context.Background(), &sql.TxOptions{
		Isolation: isoLevel,
	})
	if err != nil {
		stderr(err, 1)
		os.Exit(1)
	}

	trx := &autoTransaction{
		tx: xtx,
	}
	fn(trx)
	if tv := os.Getenv("DBMSESS__TEST_MODE"); tv == "1" {
		stdout(os.Getenv("DBMSESS__CURRENT_NAME"), "success... rolling back (TEST MODE)")
		stderr(xtx.Rollback(), 1)
		return
	}

	if _, err := xtx.Exec(os.Getenv("DBMSESS__INSERT_QUERY"), os.Getenv("DBMSESS__UUID"), os.Getenv("DBMSESS__AUTHOR"), os.Getenv("DBMSESS__CREATED")); err != nil {
		stdout(os.Getenv("DBMSESS__CURRENT_NAME"), "insert migration ID error", err.Error())
		stderr(xtx.Rollback(), 1)
		os.Exit(1)
	}

	if err := xtx.Commit(); err != nil {
		stderr(err, 1)
		os.Exit(1)
	}
}

func open(driverName, dataSourceName string) (*sql.DB, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}
	if driverCtx, ok := driveri.(driver.DriverContext); ok {
		connector, err := driverCtx.OpenConnector(dataSourceName)
		if err != nil {
			return nil, err
		}
		return sql.OpenDB(connector), nil
	}

	return sql.OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
}

func init() {
	driversMu.Lock()
	defer driversMu.Unlock()
	drivers = make(map[string]driver.Driver)
	drivers["mysql"] = &mysql.MySQLDriver{}
	drivers["sqlite3"] = &sqlite3.SQLiteDriver{}
}

func EnvConnect() (*sqlx.DB, error) {
	driverName := os.Getenv("DBMSESS__DRIVER_NAME")
	dataSourceName := os.Getenv("DBMSESS__DATA_SOURCE_NAME")
	rawdb, err := open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	xdb := sqlx.NewDb(rawdb, driverName)
	return xdb, nil
}
