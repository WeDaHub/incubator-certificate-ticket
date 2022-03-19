package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"qcloud-tools/core/config"
	"time"
)

func init() {
	var db = config.QcloudTool.Db
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", db.User, db.Password, db.Host, db.Port, db.Database)

	QcloudToolDb = Conn{
		Dsn: dsn,
	}

	QcloudToolDb.Db, err = sql.Open("mysql", QcloudToolDb.Dsn)
	if err != nil {
		fmt.Printf("failed to connect database : %v\n", err)
		os.Exit(1)
	}

	QcloudToolDb.Db.SetConnMaxLifetime(time.Duration(db.ConnMaxLifetime) * time.Second)
	QcloudToolDb.Db.SetMaxIdleConns(db.MaxIdleConn)

	err = QcloudToolDb.Db.Ping()
	if err != nil {
		fmt.Printf("failed to ping database : %v\n", err)
		os.Exit(1)
	}
}

func CloseDb(conn Conn) {
	conn.Db.Close()
}
