package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Conn struct {
	Dsn string
	Db  *sql.DB
}

const tryAgainError = "try connecting again"

const maxBadConnRetries = 2

var QcloudToolDb Conn

func (conn Conn) Delete(sqlStr string, args ...interface{}) (affected int64, err error) {
	stmt, err := conn.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if nil != err {
		fmt.Println("failed to exec query", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	return
}

func (conn Conn) Update(sqlStr string, args ...interface{}) (affected int64, err error) {

	stmt, err := conn.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if nil != err {
		fmt.Println("failed to exec query", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	return
}

func (conn Conn) Insert(sqlStr string, args ...interface{}) (lastInsertId int64, err error) {

	stmt, err := conn.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if nil != err {
		fmt.Println("failed to exec query", err)
		return 0, err
	}

	lastInsertId, err = result.LastInsertId()
	return
}

func (conn Conn) Query(sqlStr string, args ...interface{}) (rows *sql.Rows, err error) {

	stmt, err := conn.Prepare(sqlStr)
	if nil != err {
		fmt.Println("failed to prepare query:", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(args...)
	if nil != err {
		fmt.Println("failed to query data:", err)
		return
	}

	return
}

func (conn Conn) Prepare(sql string) (stmt *sql.Stmt, err error) {

	for i := 0; i < maxBadConnRetries; i++ {
		stmt, err = QcloudToolDb.Db.Prepare(sql)
		if err == nil || !strings.Contains(err.Error(), tryAgainError) {
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
	}

	return
}
