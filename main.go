package main

import (
	"fmt"
	"net/http"
	"todo-api/controllers"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", controllers.GetTodos)
	mux.HandleFunc("POST /todos", controllers.CreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", controllers.DeleteTodo)
	mux.HandleFunc("PATCH /todos/{id}", controllers.PatchTodo)
	fmt.Println("Сервер запущен на порту: 8080")
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Ошибка при запуска: ", err)
	}


}