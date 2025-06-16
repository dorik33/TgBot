package bot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) ticker() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			subs := b.supbrepo.GetAllSubs()
			fmt.Println("mailing is start")
			for _, sub := range subs {
				crypto, _ := b.apiClient.GetInfo(sub.Token)

				price := crypto.PriceUSD
				msg := tgbotapi.NewMessage(sub.UserID, fmt.Sprintf("Цена на токен %s: %f", sub.Token, price))
				b.botAPI.Send(msg)
			}
		}
	}
}
