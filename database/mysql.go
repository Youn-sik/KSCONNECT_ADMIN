package database

import (
	"database/sql"
	"log"

	"github.com/Youn-sik/KSCONNECT_ADMIN/setting"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlConnection() *sql.DB {
	config, err := setting.LoadConfigSettingJSON()

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db, err := sql.Open("mysql", config.Mysql.User+":"+config.Mysql.Password+"@tcp("+config.Mysql.Host+":"+config.Mysql.Port+")/"+config.Mysql.Database)
	if err != nil {
		panic(err)
	}
	return db
}
