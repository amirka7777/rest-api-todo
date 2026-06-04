package repository

import (
	"todo-api/database"
	"todo-api/models"
)

func GetAllTodos() ([]models.Todo, error) {

	rows, err := database.DB.Query("SELECT id, title, completed, created_at FROM todos")
	if err != nil {
		return nil, err
	}

	totalTodos := []models.Todo{}
	for rows.Next() {
		var tmp models.Todo
		err = rows.Scan(&tmp.ID, &tmp.Title, &tmp.Completed, &tmp.CreatedAt)
		if err != nil {
			return nil, err
		}
		totalTodos = append(totalTodos, tmp)
	}

	return totalTodos, nil

}

func GetTodoByID(id int) (models.Todo, error) {
	var tmp models.Todo
	err := database.DB.QueryRow("SELECT id, title, completed, created_at FROM todos WHERE id = ?", id).Scan(&tmp.ID, &tmp.Title, &tmp.Completed, &tmp.CreatedAt)

	return tmp, err
} 

func CreateNewtodo(title string) (int64, error) {

	result, err := database.DB.Exec("INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		return 0, err
	}
	// возвращаем id задачи
	return result.LastInsertId()

}

func DeleteTodoFromBD(id int) (int64, error) {
	result, err := database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}


func CustomUpdate(query string, args ...interface{}) (int64, error) {

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}