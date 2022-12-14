package goorm

import (
	"database/sql"
	"go-orm/log"
	"go-orm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, datasource string) (e *Engine, err error) {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Error(err)
		return
	}
	// send a ping to make sure the connection is alive
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	return &Engine{db: db}, nil
}
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("fail to close database...", err)
	}
	log.Info("database close success...")
}
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
