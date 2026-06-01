package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo-api/database"
	"todo-api/models"
)


func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rows, err := database.DB.Query("SELECT id, title, completed, created_at FROM todos")
	if err != nil {
		http.Error(w, "Ошибка при чтении из базы данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Todo = []models.Todo{}
	for rows.Next() {
		var t models.Todo
		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt)
		if err != nil {
			http.Error(w, "Ошибка при разборе данных из базы данных", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, t)

	}

	json.NewEncoder(w).Encode(tasks)

}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo models.Todo

	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTodo.CreatedAt = time.Now()

	var result sql.Result
	result, err = database.DB.Exec("INSERT INTO todos (title, created_at) VALUES (?, ?)", newTodo.Title, newTodo.CreatedAt)
	if err != nil {
		http.Error(w, "Ошибка при вставке в базу данных", http.StatusInternalServerError)
		return
	}
	var newID int64
	newID, err = result.LastInsertId()
	if err != nil {
		http.Error(w, "Не удалось получить ID новой задачи", http.StatusInternalServerError)
		return
	}

	newTodo.ID = int(newID)
 
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

	var result sql.Result
	result, err = database.DB.Exec("DELETE FROM todos WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, "Ошибка при удалении из базы данных", http.StatusInternalServerError)
		return
	}

	var rowsAfterDelete int64
	rowsAfterDelete, err = result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка при проверки результата", http.StatusInternalServerError)
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
		Title *string `json:"title"`
		Completed *bool `json:"completed"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	query := "UPDATE todos SET "
	var args []interface{}

	if input.Title != nil{
		query += "title = ?, "
		args = append(args, *input.Title)
	}
	if input.Completed != nil {
		query += "completed = ?, "
		args = append(args, *input.Completed)
	}

	if len(args) == 0 {
		http.Error(w, "Ошибка при передачи данных с клиента", http.StatusBadRequest)
		return
	}
	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, idInt)

	var result sql.Result
	result, err = database.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Ошибка при обновлении базы данных", http.StatusInternalServerError)
		return
	}
	
	var rowsAfterUpdate int64
	rowsAfterUpdate, err = result.RowsAffected()
	if err != nil {
		http.Error(w, "Ошибка при проверки результата", http.StatusInternalServerError)
		return
	}

	if rowsAfterUpdate == 0 {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	var updateTask models.Todo
	err = database.DB.QueryRow("SELECT id, title, completed, created_at FROM todos WHERE id = ?", idInt).Scan(&updateTask.ID, &updateTask.Title, &updateTask.Completed, &updateTask.CreatedAt)

	if err != nil {
		http.Error(w, "Ошибка при получении обновленных данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateTask)

	
	
}