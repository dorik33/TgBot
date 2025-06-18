package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/bot"
	"github.com/dorik33/TgBot/internal/config"
	redissrorage "github.com/dorik33/TgBot/internal/repository/redis_srorage"
	"github.com/dorik33/TgBot/internal/repository/store"

	subscriptionservice "github.com/dorik33/TgBot/internal/service/Subscriptionservice"
	"github.com/dorik33/TgBot/internal/service/cryptoservice"
	"github.com/dorik33/TgBot/internal/service/walletservice"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := config.Load()
	fmt.Println(cfg)

	client := api.NewAPIClient(cfg.ApiKey, cfg.TimeOut)

	store, err := store.NewConnection(*cfg)
	if err != nil {
		log.Printf("Не удалось подключиться к базе данных: %v", err)
		os.Exit(1)
	}
	defer store.Close()

	rdb, err := redissrorage.NewClient(context.Background(), *cfg)
	if err != nil {
		log.Printf("Не удалось подключиться к redis: %v", err)
		os.Exit(1)
	}
	redis := redissrorage.NewRedis(rdb, cfg.RedisConfig.TTL)
	defer redis.Close()

	botAPI, _ := tgbotapi.NewBotAPI(cfg.BotKey)

	cryptoService := cryptoservice.NewCryptoService(client, redis)
	subService := subscriptionservice.NewSubscriptionService(store.SubscriptionRepository, cryptoService)
	walletService := walletservice.NewWalletService(store.WalletRepository, cryptoService)

	bot := bot.NewBot(botAPI, subService, walletService, cryptoService)

	go bot.Ticker()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}
	log.Printf("Update channel initialized")

	logMiddleware := func(next func(tgbotapi.Update)) func(tgbotapi.Update) {
		return func(update tgbotapi.Update) {
			if update.Message != nil {
				log.Printf("Received message from chat %d at %s: %s",
					update.Message.Chat.ID,
					time.Now().Format(time.RFC3339),
					update.Message.Text)
			}
			next(update)
		}
	}

	wrappedHandler := logMiddleware(func(update tgbotapi.Update) {
		bot.HandleUpdate(update)
	})

	for update := range updates {
		wrappedHandler(update)
	}

}
