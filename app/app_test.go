package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"jhpark.sinsiway.com/bootstrap-oauth/model"
)

func TestTodos(t *testing.T) {

	assert := assert.New(t)
	ts := httptest.NewServer(MakeHandler("./test.db"))
	defer ts.Close()
	//
	// POST /todos
	//
	// test todos1
	//
	res, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todos1"}})
	assert.NoError(err)

	assert.Equal(http.StatusOK, res.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todos1")

	id1 := todo.Id

	//
	// add test todos2
	//
	res, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todos2"}})
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todos2")

	id2 := todo.Id

	//
	// GET /todos
	//
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	var todoList []*model.Todo

	err = json.NewDecoder(res.Body).Decode(&todoList)
	assert.NoError(err)
	assert.Equal(2, len(todoList))

	for _, t := range todoList {
		if t.Id == id1 {
			assert.Equal("Test todos1", t.Name)
		} else if t.Id == id2 {
			assert.Equal("Test todos2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	//
	// Get /complete-todo/
	//
	res, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	//
	// GET /todos
	//
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	todoList = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todoList)
	assert.NoError(err)

	for _, t := range todoList {
		if t.Id == id1 {
			assert.True(t.Completed)
		}
	}

	//
	// DELETE /todos
	//
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos?id="+strconv.Itoa(id1), nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	//
	// GET /todos
	//
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	todoList = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todoList)
	assert.NoError(err)

	assert.Equal(1, len(todoList))
}
