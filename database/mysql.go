package database

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBParameter struct {
	Host         string
	Port         int
	Username     string
	Password     string
	Dbname       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

func CreateMysqlDBxWithDBParameter(p DBParameter) (*sqlx.DB, error) {
	strPort := strconv.Itoa(p.Port)
	return CreateMysqlDBx(p.Host, strPort, p.Username, p.Password, p.Dbname, p.Charset, p.MaxIdleConns, p.MaxOpenConns)
}

func CreateMysqlDBx(host, port, username, password, dbname, charset string, idle, open int) (*sqlx.DB, error) {
	db, err := CreateMysqlDB(host, port, username, password, dbname, charset, idle, open)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "mysql")
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return dbx, nil
}

func CreateMysqlDBWithDBParameter(p DBParameter) (*sql.DB, error) {
	strPort := strconv.Itoa(p.Port)
	return CreateMysqlDB(p.Host, strPort, p.Username, p.Password, p.Dbname, p.Charset, p.MaxIdleConns, p.MaxOpenConns)
}

func CreateMysqlDB(host, port, username, password, dbname, charset string, idle, open int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, port, dbname, charset)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("can't not connect to db: host = %s, port = %s, username = %s, dbname = %s, err = %v", host, port, username, dbname, err))
	}
	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(open)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
