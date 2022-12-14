package model

import (
	"time"
)

type Todo struct {
	Id        int       `json:"id" gorm:"column:id;primary key;autoincrement"`
	Name      string    `json:"name" gorm:"column:name"`
	Completed bool      `json:"completed" gorm:"column:completed"`
	CreateAt  time.Time `json:"create_at" gorm:"column:create_at"`
}

type DbHandler interface {
	GetTodos() []*Todo
	AddTodo(name string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool
	Close()
}

func NewDbHandler(filepath string) DbHandler {
	return newSqlite3Handler(filepath)
}

/*
func GetTodos() []*Todo {
	return handler.GetTodos()
}

func AddTodo(name string) *Todo {
	return handler.AddTodo(name)
}

func RemoveTodo(id int) bool {
	return handler.RemoveTodo(id)
}

func CompleteTodo(id int, complete bool) bool {
	return handler.CompleteTodo(id, complete)
}
*/
