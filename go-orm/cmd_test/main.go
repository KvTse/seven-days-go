package main

import (
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	goorm "go-orm"
	"go-orm/log"
)

const (
	driverName     = "mysql"
	dataSourceName = "root:root@(localhost:3306)/test"
)

func main() {
	e, err := goorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return
	}
	s := e.NewSession()
	result, err := s.Raw("insert into user (`name`) values (?),(?)", "Tony", "Sandy").Exec()
	affected, _ := result.RowsAffected()
	log.Infof("change %d row ", affected)
}
