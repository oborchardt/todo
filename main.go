package main

import (
	"log"
	"net/http"
	"todo/controllers"
	"todo/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", controllers.LoginUser)
	mux.HandleFunc("GET /users", controllers.GetUsers)
	mux.HandleFunc("POST /users", controllers.CreateUser)
	mux.Handle("GET /todos", middlewares.Authenticate(http.HandlerFunc(controllers.GetTodos)))
	mux.Handle("POST /todos", middlewares.Authenticate(http.HandlerFunc(controllers.CreateTodo)))
	mux.Handle("DELETE /todos", middlewares.Authenticate(http.HandlerFunc(controllers.DeleteTodo)))
	// according to https://stackoverflow.com/questions/28459418/use-of-put-vs-patch-methods-in-rest-api-real-life-scenarios
	mux.Handle("PATCH /todos", middlewares.Authenticate(http.HandlerFunc(controllers.UpdateTodo)))

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
