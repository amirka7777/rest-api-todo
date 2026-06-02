package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite3", "todos.db")
	if err != nil {
		log.Fatalf("Не удалось подключить драйвера базы данных: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("База данных не доступна: %v", err)
	}

	createTableSQLite := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTableSQLite)
	if err != nil {
		log.Fatalf("Ошибка при создании базы данных: %v\n", err)
	}

	log.Println("База данных была проверена и успешно подключена")
}