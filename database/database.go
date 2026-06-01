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

	log.Println("База данных была проверена и успешно подключена")
}