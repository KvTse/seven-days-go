package session

import (
	"fmt"
	"go-orm/log"
	"go-orm/schema"
	"reflect"
	"strings"
)

// Model 用于给refTable 赋值 refTable->schema 属于是表的信息
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 返回refTable的值
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("model is not set")
	}
	return s.refTable
}

// CreateTable 根据当前Session中存在的内容,创建表
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprint("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("Create table %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable 删除表
func (s *Session) DropTable() error {
	table := s.RefTable()
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.Name)).Exec()
	return err
}

// HasTable 判断是否存在表
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var temp string
	_ = row.Scan(&temp)
	return temp == s.RefTable().Name
}
