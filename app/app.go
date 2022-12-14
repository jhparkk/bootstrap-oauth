package app

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"jhpark.sinsiway.com/bootstrap-oauth/model"
)

var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db model.DbHandler
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	//var list []*model.Todo
	//for _, v := range todoMap {
	//	log.Println("v:", v)
	//	list = append(list, v)
	//}
	list := a.db.GetTodos()
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	todo := a.db.AddTodo(name)

	//id := len(todoMap) + 1
	//todo := &Todo{id, name, false, time.Now()}
	//todoMap[id] = todo

	rd.JSON(w, http.StatusOK, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, _ := strconv.Atoi(vars["id"])

	id, _ := strconv.Atoi(r.FormValue("id"))
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
	//if _, ok := todoMap[id]; ok {
	//	delete(todoMap, id)
	//	rd.JSON(w, http.StatusOK, Success{true})
	//	log.Println("deleted id:", id)
	//} else {
	//	rd.JSON(w, http.StatusOK, Success{false})
	//	log.Println("not found id:", id)
	//}

}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"

	ok := a.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}

	//if todo, ok := todoMap[id]; ok {
	//	todo.Completed = complete
	//	rd.JSON(w, http.StatusOK, Success{true})
	//} else {
	//	rd.JSON(w, http.StatusOK, Success{false})
	//}
}

//func addTestTodo() {
//	todoMap[1] = &Todo{1, "Init Data1", false, time.Now()}
//	todoMap[2] = &Todo{2, "Init Data2", true, time.Now()}
//	todoMap[3] = &Todo{3, "Init Data3", false, time.Now()}
//}

func (a *AppHandler) Close() {
	a.db.Close()
}

func MakeHandler(filepath string) *AppHandler {
	//todoMap = make(map[int]*Todo)
	//addTestTodo()

	r := mux.NewRouter()
	app := &AppHandler{
		Handler: r,
		db:      model.NewDbHandler(filepath),
	}
	r.HandleFunc("/todos", app.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", app.addTodoListHandler).Methods("POST")
	r.HandleFunc("/todos", app.removeTodoHandler).Methods("DELETE")
	//r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", app.completeTodoHandler).Methods("GET")

	r.HandleFunc("/", app.indexHandler)
	return app
}
