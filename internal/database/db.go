package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/config"
)

func NewConnection(cfg config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.Database.DBCon)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных")
	}
	fmt.Println(cfg.Database.DBCon)
	if err = db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения")
	}
	log.Println("База данных успешно запущена")

	return db
}
