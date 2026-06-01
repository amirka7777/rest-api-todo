package main

import (
	"fmt"
	"log"
	"net/http"
	"todo-api/controllers"
	"todo-api/database"
)

func main() {

	database.InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", controllers.GetTodos)
	mux.HandleFunc("POST /todos", controllers.CreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", controllers.DeleteTodo)
	mux.HandleFunc("PATCH /todos/{id}", controllers.PatchTodo)
	log.Println("Сервер запущен на порту: 8080")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Ошибка при запуска: ", err)
	}


}