package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/bot"
	"github.com/dorik33/TgBot/internal/config"
	"github.com/dorik33/TgBot/internal/database"
)

func main() {
	cfg := config.Load()
	os.Setenv("DATABASE_USER", cfg.Database.User)
	os.Setenv("DATABASE_PASSWORD", cfg.Database.Password)
	os.Setenv("DATABASE_NAME", cfg.Database.DBName)

	client := api.NewAPIClient(cfg.TimeOut)
	data, err := client.GetInfo("bitcoin")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(data.PriceUSD)

	db := database.NewConnection(*cfg)
	subRepo := database.NewSubscriptionRepository(db)

	subs, err := subRepo.GetSubcriptions(914333594)
	if err != nil {
		log.Fatalf("не работает бд  ", err)
	}
	log.Println(subs)

	bot.StartBot(cfg.BotKey, client, subRepo)

	//TODO Добавить readme, добавить кеширование, возможно добавить Make
}
