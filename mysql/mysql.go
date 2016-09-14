package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	MYSQL_CONN_CAP = 5
	err            error
	helperLock     *sync.Mutex = &sync.Mutex{}
	dbHelper       DBHelper
)

func init() {
	fmt.Println(" mysql init...")
}

type DBHelper interface {
	Insert(a, b string, ch chan int)
}

type dbStruct struct {
	db *sql.DB
}

func GetDBHelper() DBHelper {
	helperLock.Lock()
	defer helperLock.Unlock()
	if dbHelper != nil {
		return dbHelper
	} else {
		db, err := sql.Open("mysql", "test:test@tcp(127.0.0.1:3306)/test?charset=utf8")
		if err != nil {
			panic("open db err")
		}
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(5)
		dbHelper = &dbStruct{db}
		return dbHelper
	}

}

func (this *dbStruct) Insert(a, b string, ch chan int) {
	tx, err := this.db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := tx.Prepare("INSERT INTO user(username, password) VALUES(?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 1000; i++ {
		stmt.Exec(a+strconv.Itoa(i), b)
	}
	tx.Commit()
	ch <- 1
}
