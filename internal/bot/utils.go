package bot

import (
	"fmt"

	"github.com/dorik33/TgBot/internal/models"
)

func parseSubs(subs []models.Subscription) string {
	res := "Ваши подписки:\n"
	for i, sub := range subs {
		str := fmt.Sprintf("Подписка %d на токен %s\n", i+1, sub.Token)
		res += str
	}
	return res
}
