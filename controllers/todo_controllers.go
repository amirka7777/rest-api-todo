package controllers

// w http.ResponseWriter - инструмент для отправки на клиент
//  r *http.Request - указатель на сам запрос

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	for i, task := range Tasks {
		if task.ID == idInt {
			Tasks = append(Tasks[:i], Tasks[i+1:]... )
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Задача не найдена", http.StatusNotFound)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest) 
		return
	}

	var updateTask models.Todo
	err = json.NewDecoder(r.Body).Decode(&updateTask)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	for i, task := range Tasks {
		if task.ID == idInt {
			Tasks[i].Title = updateTask.Title
			Tasks[i].Completed = updateTask.Completed
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Tasks[i])
			return
		}
	}

	http.Error(w, "Задача не найдена", http.StatusNotFound)
	
}