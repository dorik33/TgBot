package bot

import (
	"log"
	"strings"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartBot(key string, apiClient *api.APIClient, rep *database.SubscriptionRepository) {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Fatal("the bot is not running")
	}
	log.Println("bot is running")
	go ticker(bot, apiClient, rep)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // Игнорируем не сообщения
			continue
		}

		// Выводим текст полученного сообщения
		log.Printf("[%s] %s, %d\n", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

		// Ответ на команду /start
		if update.Message.Text == "/start" {
			StartMessage(bot, update)
			log.Println(update.Message.Chat.ID)
		}

		// Ответ на команду /help
		if update.Message.Text == "/help" {
			HelpMessage(bot, update)
		}
		if strings.Split(update.Message.Text, " ")[0] == "/price" {
			PriceMessage(bot, update, apiClient)
		}
		if strings.Split(update.Message.Text, " ")[0] == "/sub" {
			AddSubscriptionMessage(bot, update, apiClient, rep)
		}
		if strings.Split(update.Message.Text, " ")[0] == "/subs" {
			GetSubcriptionsMessage(bot, update, apiClient, rep)
		}
		if strings.Split(update.Message.Text, " ")[0] == "/delete_sub" {
			DeleteSubcriptionsMessage(bot, update, apiClient, rep)
		}
	}
}
