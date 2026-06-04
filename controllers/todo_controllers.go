package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/models"
	"todo-api/repository"
)


func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	tasks, err := repository.GetAllTodos()
	if err != nil {
		http.Error(w, "Ошибка при получении данных из базы данных", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	
	var input struct {
		Title string `json:"title"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}
	var newID int64
	newID, err = repository.CreateNewtodo(input.Title)
	if err != nil {
		http.Error(w, "Ошибка при создании задачи", http.StatusInternalServerError)
		return
	}

	var newTodo models.Todo
	newTodo, err = repository.GetTodoByID(int(newID))
	if err != nil {
		http.Error(w, "Ошибка при получении задачи из Базы данных", http.StatusInternalServerError)
		return
	}

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

	var rowsAfterDelete int64
	rowsAfterDelete, err = repository.DeleteTodoFromBD(idInt)
	if err != nil {
		http.Error(w, "Ошибка при удалении из базы данных", http.StatusInternalServerError)
		return
	}

	if rowsAfterDelete == 0 {
		http.Error(w, "Задача не была найдена", http.StatusNotFound)
		return
	}


	w.WriteHeader(http.StatusNoContent)
}

func PatchTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	query := "UPDATE todos SET "
	var args []interface{}

	if input.Title != nil {
		query += "title = ?, "
		args = append(args, *input.Title)
	}
	if input.Completed != nil {
		query += "completed = ?, "
		args = append(args, *input.Completed)
	}

	if len(args) == 0 {
		http.Error(w, "Нет данных для обновления", http.StatusBadRequest)
		return
	}

	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, idInt)

	rowsAffected, err := repository.CustomUpdate(query, args...)
	if err != nil {
		http.Error(w, "Ошибка при обновлении базы данных", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	updatedTask, err := repository.GetTodoByID(idInt)
	if err != nil {
		http.Error(w, "Ошибка при получении обновленных данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}