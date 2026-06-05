package repository

import (
	"todo-api/database"
	"todo-api/models"
)

func GetAllTodos(search string, completed *bool) ([]models.Todo, error) {

	query := "SELECT id, title, completed, created_at FROM todos WHERE 1=1"
	var args []interface{}

	if search != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+search+"%")
	}

	if completed != nil {
		query += " AND completed = ?"
		args = append(args, *completed)
	}

	query += " ORDER BY created_at DESC"


	rows, err := database.DB.Query(query, args...)
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