package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqlite3Handler struct {
	db *gorm.DB
}

func (s *sqlite3Handler) GetTodos() []*Todo {
	var todos []*Todo
	s.db.Table("todos").Find(&todos)
	return todos
}

func (s *sqlite3Handler) AddTodo(name string) *Todo {
	todo := Todo{Name: name, Completed: false, CreateAt: time.Now()}
	tx := s.db.Table("todos").Create(&todo)
	if tx.Error != nil {
		panic(tx.Error)
	}
	return &todo
}

func (s *sqlite3Handler) RemoveTodo(id int) bool {
	tx := s.db.Table("todos").Delete(&Todo{}, id)
	if tx.Error != nil {
		return false
	}
	return true
}

func (s *sqlite3Handler) CompleteTodo(id int, complete bool) bool {
	var todo Todo
	tx := s.db.Table("todos").Where(id).Find(&todo)
	if tx.Error != nil {
		return false
	}

	todo.Completed = complete
	tx = tx.Save(&todo)
	if tx.Error != nil {
		return false
	}

	return true
}

func (s *sqlite3Handler) Close() {
	db, _ := s.db.DB()
	db.Close()
}

func newSqlite3Handler(filepath string) DbHandler {
	db, err := gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	user := &Todo{}
	err = db.AutoMigrate(user)
	if err != nil {
		panic(err)
	}

	return &sqlite3Handler{db: db}
}
