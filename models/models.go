package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lhw0828/go-gin-example/pkg/setting"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         uint `json:"id" gorm:"primary_key"`
	CreatedOn  int  `json:"created_on"`
	ModifiedOn int  `json:"modified_on"`
}

func init() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	sec, err := setting.Cfg.GetSection("database")

	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	db, err = gorm.Open(dbType, user+":"+password+"@tcp("+host+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}
	db.SingularTable(true)
	db.LogMode(true)
	if db != nil {
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
	} else {
		log.Fatal("Failed to initialize database connection")
	}

}

func CloseDB() {
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)
}
