package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/config"
)

func NewConnection(cfg config.Config) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных")
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения")
	}
	log.Println("База данных успешно запущена")

	return db
}
