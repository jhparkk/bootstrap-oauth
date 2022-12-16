package model

import (
	"log"
	"testing"
)

func TestSqlite3Handler(t *testing.T) {
	dbhandler := newSqlite3Handler("./test.db")
	sessId := "sess_id:test"
	dbhandler.AddTodo(sessId, "test1")
	dbhandler.AddTodo(sessId, "test2")
	dbhandler.AddTodo(sessId, "test3")
	dbhandler.AddTodo("none", "none1")
	dbhandler.AddTodo("none", "none2")

	todos := dbhandler.GetTodos(sessId)
	for idx, v := range todos {
		log.Println(idx, " : ", v)
	}
}
