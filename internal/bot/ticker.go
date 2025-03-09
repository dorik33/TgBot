package bot

import (
	"fmt"
	"time"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ticker(bot *tgbotapi.BotAPI, client *api.APIClient, repo *database.SubscriptionRepository) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			subs := repo.GetAllSubs()
			fmt.Println("mailing is start")
			for _, sub := range subs {
				crypto, _ := client.GetInfo(sub.Token)

				price := parsePrice(crypto.PriceUSD)
				msg := tgbotapi.NewMessage(sub.UserID, fmt.Sprintf("Цена на токен %s: %f", sub.Token, price))
				bot.Send(msg)
			}
		}
	}
}
