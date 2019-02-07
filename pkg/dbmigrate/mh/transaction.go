package mh

import (
	"context"
	"database/sql"
	"os"

	"github.com/jmoiron/sqlx"
)

type ExitTx interface {
	BindNamed(query string, arg interface{}) (string, []interface{})
	DriverName() string
	Exec(query string, args ...interface{}) sql.Result
	ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	Get(dest interface{}, query string, args ...interface{})
	//Rebind(query string) string
	//Unsafe() Transaction
	//NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	//NamedExec(query string, arg interface{}) (sql.Result, error)
	//Select(dest interface{}, query string, args ...interface{}) error
	//Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	//QueryRowx(query string, args ...interface{}) *sqlx.Row
	//MustExec(query string, args ...interface{}) sql.Result
	//Preparex(query string) (*sqlx.Stmt, error)
	//Stmtx(stmt interface{}) *Stmt
	//NamedStmt(stmt *NamedStmt) *NamedStmt
	//PrepareNamed(query string) (*NamedStmt, error)
}

type autoTransaction struct {
	tx *sqlx.Tx
}

func (a *autoTransaction) BindNamed(query string, arg interface{}) (string, []interface{}) {
	x, y, err := a.tx.BindNamed(query, arg)
	if err != nil {
		stderr(err, 2)
		stderr(a.tx.Rollback(), 2)
		os.Exit(1)
	}
	return x, y
}

func (a *autoTransaction) DriverName() string {
	return a.tx.DriverName()
}

func (a *autoTransaction) Exec(query string, args ...interface{}) sql.Result {
	result, err := a.tx.Exec(query, args...)
	if err != nil {
		stderr(err, 2)
		stderr(a.tx.Rollback(), 2)
		os.Exit(1)
	}
	return result
}

func (a *autoTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
	res, err := a.tx.ExecContext(ctx, query, args...)
	if err != nil {
		stderr(err, 2)
		stderr(a.tx.Rollback(), 2)
		os.Exit(1)
	}
	return res
}

func (a *autoTransaction) Get(dest interface{}, query string, args ...interface{}) {
	err := a.tx.Get(dest, query, args...)
	if err != nil {
		stderr(err, 2)
		stderr(a.tx.Rollback(), 2)
		os.Exit(1)
	}
}
