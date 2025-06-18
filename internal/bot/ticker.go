package bot

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) Ticker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			subs, err := b.subscriptionService.GetAllSubscriptions()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("mailing is start")
			for _, sub := range subs {
				crypto, _ := b.cryptoService.GetCryptoPrice(sub.Token)

				price := crypto.PriceUSD
				msg := tgbotapi.NewMessage(sub.UserID, fmt.Sprintf("Цена на токен %s: %f", sub.Token, price))
				b.botAPI.Send(msg)
			}
		}
	}
}
