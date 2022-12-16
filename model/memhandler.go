package model

import "time"

type memoryHandler struct {
	todoMap map[int]*Todo
}

func (m *memoryHandler) GetTodos(sessId string) []*Todo {
	var list []*Todo
	for _, v := range m.todoMap {
		//log.Println("v:", v)
		if v.SessionId == sessId {
			list = append(list, v)
		}
	}
	return list
}

func (m *memoryHandler) AddTodo(sessId, name string) *Todo {
	id := len(m.todoMap) + 1
	todo := &Todo{id, sessId, name, false, time.Now()}
	m.todoMap[id] = todo

	return todo
}

func (m *memoryHandler) RemoveTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok {
		delete(m.todoMap, id)
		return true
	}

	return false
}

func (m *memoryHandler) CompleteTodo(id int, complete bool) bool {
	if todo, ok := m.todoMap[id]; ok {
		todo.Completed = complete
		return true
	}

	return false
}

func (m *memoryHandler) Close() {

}

func newMemoryHandler() DbHandler {
	m := &memoryHandler{}
	m.todoMap = make(map[int]*Todo)
	return m
}
