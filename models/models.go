package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lhw0828/go-gin-example/pkg/setting"
	"log"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         uint `json:"id" gorm:"primary_key"`
	CreatedOn  int  `json:"created_on"`
	ModifiedOn int  `json:"modified_on"`
	DeletedOn  int  `json:"deleted_on"`
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
		// 增加回调函数，在插入和更新时自动更新时间
		db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
		db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
		// 增加回调函数，判断是软删除还是硬删除
		db.Callback().Delete().Replace("gorm:delete", deleteCallback)
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

// updateTimeStampForUpdateCallback will set `CreatedOn “ModifiedOn` when create callback is triggered with option `gorm:update_time_stamp=true`
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				err := createTimeField.Set(nowTime)
				if err != nil {
					return
				}
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				err := modifyTimeField.Set(nowTime)
				if err != nil {
					return
				}
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when update callback is triggered with option `gorm:update_time_stamp=true`
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		err := scope.SetColumn("ModifiedOn", time.Now().Unix())
		if err != nil {
			return
		}
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()), addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
