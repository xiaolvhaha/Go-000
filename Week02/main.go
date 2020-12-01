package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type user struct {
	ID   uint `db:"id"`
	Name string `db:"name"`
}

var Db  *sqlx.DB

func init(){
	database, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	Db = database

}


func main() {
	usr, err := service(10000)

	if err != nil {
		if errors.Is(errors.Cause(err), sql.ErrNoRows) {
			fmt.Println("not found")
			return
		}

		fmt.Println("sth go wrong")
		return
	}
	fmt.Println("200", usr.Name)
	defer Db.Close()
}


func service(id uint) (*user, error) {
	return dao(id)
}

func dao(id uint) (*user, error) {
	var usr user
	err := Db.QueryRow("select id, name from user where id = ?", id).Scan(&usr.ID, &usr.Name)
	if err == sql.ErrNoRows {

		return nil, errors.Wrapf(err, "user can not found by id:%d", id)
	}
	if err != nil {
		return nil, errors.Wrap(err, "err")
	}
	return &usr, nil
}
