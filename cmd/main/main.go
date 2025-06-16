package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/bot"
	"github.com/dorik33/TgBot/internal/config"
	"github.com/dorik33/TgBot/internal/database"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)

	client := api.NewAPIClient(cfg.ApiKey, cfg.TimeOut)

	db := database.NewConnection(*cfg)
	subRepo := database.NewSubscriptionRepository(db)
	walletrepo := database.NewWalletRepository(db)

	subs, err := subRepo.GetSubcriptions(914333594)
	if err != nil {
		log.Fatalf("не работает бд: %v  ", err)
	}
	log.Println(subs)

	b, err := bot.NewBot(cfg.BotKey, client, subRepo, walletrepo)
	if err != nil {
		log.Fatalf("error while start bot: %v", err)
	}
	b.Start()

	//TODO Добавить readme, добавить кеширование, возможно добавить Make
}
