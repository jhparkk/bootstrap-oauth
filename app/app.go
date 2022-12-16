package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"jhpark.sinsiway.com/bootstrap-oauth/model"

	"github.com/google/uuid"
)

func createSessionKey() string {
	id := uuid.New()
	sessKey := id.String()
	return "jhpark_sesskey:" + sessKey
}

var store = sessions.NewCookieStore([]byte(createSessionKey()))

var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db model.DbHandler
}

var getSessionId = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}
	id := session.Values["id"]
	if id == nil {
		return ""
	}
	return id.(string)
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
	sessId := getSessionId(r)
	list := a.db.GetTodos(sessId)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	sessId := getSessionId(r)
	todo := a.db.AddTodo(sessId, name)

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

func checkSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	log.Println("checkSignin r.URL.Path : ", r.URL.Path)

	// if request URL is /signin.html, then next()
	if strings.Contains(r.URL.Path, "/signin") || strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	sessId := getSessionId(r)
	log.Println("checkSignin sessId : ", sessId)
	// if user already signed in
	if sessId != "" {
		next(w, r)
		return
	}

	// if not user signed in then redirect signin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

func MakeHandler(filepath string) *AppHandler {
	//todoMap = make(map[int]*Todo)
	//addTestTodo()

	r := mux.NewRouter()
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(checkSignin),
		negroni.NewStatic(http.Dir("public")))
	//n := negroni.Classic()  //return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))

	n.UseHandler(r)

	app := &AppHandler{
		Handler: n,
		db:      model.NewDbHandler(filepath),
	}
	r.HandleFunc("/todos", app.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", app.addTodoListHandler).Methods("POST")
	r.HandleFunc("/todos", app.removeTodoHandler).Methods("DELETE")
	//r.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", app.completeTodoHandler).Methods("GET")

	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallbackHandler) //승인된 리다이렉션 URI
	r.HandleFunc("/", app.indexHandler)

	return app
}
