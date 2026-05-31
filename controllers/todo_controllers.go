package controllers

// w http.ResponseWriter - инструмент для отправки на клиент
//  r *http.Request - указатель на сам запрос

import (
	"encoding/json"
	"net/http"
	"todo-api/models"
)

var Tasks []models.Todo = []models.Todo{
	{ID: 1, Title: "Покушать", Completed: true },
	{ID: 2, Title: "Сделать уроки", Completed: false},
}

var CurrentID int = 2

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Tasks)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo models.Todo

	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	CurrentID++
	newTodo.ID = CurrentID
	Tasks = append(Tasks, newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}