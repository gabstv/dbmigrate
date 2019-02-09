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
	"github.com/mattn/go-sqlite3"
)

var drivers map[string]driver.Driver
var driversMu sync.RWMutex

func Run(fn func(tx Mtx)) {
	driverName := os.Getenv("DBMSESS__DRIVER_NAME")
	dataSourceName := os.Getenv("DBMSESS__DATA_SOURCE_NAME")
	rawdb, err := open(driverName, dataSourceName)
	if err != nil {
		stderr(err, 1)
		os.Exit(1)
	}
	xdb := sqlx.NewDb(rawdb, driverName)
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
