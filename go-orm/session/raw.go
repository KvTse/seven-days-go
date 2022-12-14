package session

import (
	"database/sql"
	"go-orm/log"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	result, err = s.db.Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return result, nil
}

// QueryRow 查询一行数据
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

// Query 按sql查询所有
func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	rows, err = s.db.Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return
}
