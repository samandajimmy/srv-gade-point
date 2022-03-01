package database

import (
	"database/sql"
	"gade/srv-gade-point/logger"
	"strconv"
	"strings"
)

// Transaction is an interface that models the standard transaction in
// `database/sql`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// TxFn is a function that will be called with an initialized `Transaction` object.
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		logger.Make(nil, nil).Debug(err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			logger.Make(nil, nil).Fatal(p)
		}

		if err != nil {
			// something went wrong, rollback
			logger.Make(nil, nil).Debug(err)
			_ = tx.Rollback()
			return
		}

		// all good, commit
		if err := tx.Commit(); err != nil {
			logger.Make(nil, nil).Debug(err)
		}
	}()

	err = fn(tx)
	return err
}

// A PipelineStmt is a simple wrapper for creating a statement consisting of
// a query and a set of arguments to be passed to that query.
type PipelineStmt struct {
	query   string
	filters []string
	args    []interface{}
}

// NewPipelineStmt is a function to define the queries in a pipeline
func NewPipelineStmt(query string, filters []string, args ...interface{}) *PipelineStmt {
	return &PipelineStmt{query, filters, args}
}

// Exec the statement within supplied transaction. The literal string `{LAST_INS_ID}`
// will be replaced with the supplied value to make chaining `PipelineStmt` objects together
// simple.
func (ps *PipelineStmt) Exec(tx Transaction, lastInsertID int64) (sql.Result, error) {
	query := strings.Replace(ps.query, "{LAST_INS_ID}", strconv.Itoa(int(lastInsertID)), -1)
	return tx.Exec(query, ps.args...)
}

// QueryRow the statement within supplied transaction.`
func (ps *PipelineStmt) QueryRow(tx Transaction, params map[string]interface{}) *sql.Row {
	query := ps.query

	for k, v := range params {
		query = strings.Replace(query, "{"+k+"}", strconv.Itoa(int(v.(int64))), -1)
	}

	return tx.QueryRow(query, ps.args...)
}

// RunPipeline the supplied statements within the transaction. If any statement fails,
// the transaction is rolled back, and the original error is returned.
//
// The `LastInsertId` from the previous statement will be passed to `Exec`. The zero-value (0) is
// used initially.
func RunPipeline(tx Transaction, stmts ...*PipelineStmt) (sql.Result, error) {
	var res sql.Result
	var err error
	var lastInsID int64

	for _, ps := range stmts {
		res, err = ps.Exec(tx, lastInsID)
		if err != nil {
			return nil, err
		}

		lastInsID, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// RunPipelineQueryRow the supplied statements within the transaction. If any statement fails,
// the transaction is rolled back, and the original error is returned.
func RunPipelineQueryRow(tx Transaction, stmts ...*PipelineStmt) error {
	var returnedID int64
	params := map[string]interface{}{}

	for _, ps := range stmts {
		err := ps.QueryRow(tx, params).Scan(&returnedID)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if ps.filters != nil {
			params[ps.filters[0]] = returnedID
		}
	}

	return nil
}
