package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dorik33/TgBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) handlePrice(update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		b.sendMessage(update.Message.Chat.ID, "Укажите название криптовалюты: /price <имя>")
		return
	}

	crypto, err := b.cryptoService.GetCryptoPrice(args)
	if err != nil {
		log.Printf("Ошибка при запросе цены: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Не удалось получить данные 😢")
		return
	}

	price := fmt.Sprintf("💰 Цена %s: $%f", crypto.Symbol, crypto.PriceUSD)
	b.sendMessage(update.Message.Chat.ID, price)
}

func (b *Bot) handleSub(update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		b.sendMessage(update.Message.Chat.ID, "Укажите название криптовалюты: /sub <имя>")
		return
	}

	err := b.subscriptionService.Subscribe(update.Message.Chat.ID, args)
	if err != nil {
		log.Printf("Ошибка при добавлении подписки: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Не удалось добавить подписку 😢")
		return
	}

	b.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Подписка на %s успешно добавлена!", args))
}

func (b *Bot) handleListSubs(update tgbotapi.Update) {
	subs, err := b.subscriptionService.GetUserSubscriptions(update.Message.Chat.ID)
	if err != nil {
		log.Printf("Ошибка при получении подписок: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Не удалось получить подписки 😢")
		return
	}

	response := parseSubs(subs)
	b.sendMessage(update.Message.Chat.ID, response)
}

func (b *Bot) handleDeleteSub(update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		b.sendMessage(update.Message.Chat.ID, "Укажите токен: /delete_sub <имя>")
		return
	}

	err := b.subscriptionService.Unsubscribe(update.Message.Chat.ID, args)
	if err != nil {
		log.Printf("Ошибка при удалении подписки: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Не удалось удалить подписку 😢")
		return
	}

	b.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Подписка на %s успешно удалена!", args))
}

func (b *Bot) HandleAddCrypto(update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	parts := strings.Split(args, " ")

	if len(parts) < 2 {
		b.sendMessage(update.Message.Chat.ID, "Использование: /add_crypto <токен> <количество> [цена_покупки]\nПримеры:\n/add_crypto BTC 0.5\n/add_crypto ETH 2 3500")
		return
	}

	token := strings.ToUpper(parts[0])
	amountStr := parts[1]

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		b.sendMessage(update.Message.Chat.ID, "Некорректное количество. Используйте число больше 0.\nПример: 0.5 или 2")
		return
	}

	var buyPrice float64
	if len(parts) >= 3 {
		buyPrice, err = strconv.ParseFloat(parts[2], 64)
		if err != nil || buyPrice <= 0 {
			b.sendMessage(update.Message.Chat.ID, "Некорректная цена покупки. Используйте число больше 0.\nПример: 3500")
			return
		}
	}

	portfolio := &models.Portfolio{
		UserID: update.Message.Chat.ID,
		Token:  token,
		Amount: amount,
		Price:  buyPrice,
	}

	err = b.walletService.AddCryptoToWallet(portfolio)
	if err != nil {
		log.Printf("Ошибка добавления в портфолио: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Ошибка при сохранении 😢")
		return
	}

	response := fmt.Sprintf(
		"✅ Успешно добавлено:\n%s: %f по цене $%f\nОбщая стоимость покупки: $%f",
		token,
		amount,
		portfolio.Price,
		amount*portfolio.Price,
	)
	b.sendMessage(update.Message.Chat.ID, response)
}

func (b *Bot) HandleMyWallet(update tgbotapi.Update) {
	wallet, err := b.walletService.GetWallet(update.Message.Chat.ID)
	if err != nil {
		log.Printf("Ошибка получения портфолио: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Ошибка при загрузке портфолио 😢")
		return
	}

	if len(wallet) == 0 {
		b.sendMessage(update.Message.Chat.ID, "Ваше портфолио пусто.\nДобавьте активы с помощью /add_crypto")
		return
	}

	currentPrices := make(map[string]float64)
	for _, item := range wallet {
		if _, exists := currentPrices[item.Token]; !exists {
			crypto, err := b.cryptoService.GetCryptoPrice(item.Token)
			if err != nil {
				log.Printf("Ошибка получения цены для %s: %v", item.Token, err)
				currentPrices[item.Token] = 0
			} else {
				currentPrices[item.Token] = crypto.PriceUSD
			}
		}
	}

	var totalBuyCost, totalCurrentValue float64
	var builder strings.Builder
	builder.WriteString("💰 *Ваше крипто-портфолио:*\n\n")

	for _, item := range wallet {
		currentPrice := currentPrices[item.Token]
		currentValue := item.Amount * currentPrice
		buyCost := item.Amount * item.Price
		profit := currentValue - buyCost
		profitPercent := 0.0
		if buyCost != 0 {
			profitPercent = (profit / buyCost) * 100
		}

		totalBuyCost += buyCost
		totalCurrentValue += currentValue

		profitIcon := "📉"
		if profit >= 0 {
			profitIcon = "📈"
		}

		builder.WriteString(fmt.Sprintf(
			"🆔 *ID:* %d\n"+
				"🪙 *Токен:* %s\n"+
				"📦 *Количество:* %f\n"+
				"💰 *Цена покупки:* $%f\n"+
				"🏷 *Текущая цена:* $%f\n"+
				"%s *Прибыль/убыток:* $%f (%.2f%%)\n\n",
			item.ID,
			item.Token,
			item.Amount,
			item.Price,
			currentPrice,
			profitIcon,
			profit,
			profitPercent,
		))
	}

	totalProfit := totalCurrentValue - totalBuyCost
	totalProfitPercent := 0.0
	if totalBuyCost != 0 {
		totalProfitPercent = (totalProfit / totalBuyCost) * 100
	}

	totalProfitIcon := "📉"
	if totalProfit >= 0 {
		totalProfitIcon = "📈"
	}

	builder.WriteString("💹 *Итоговая статистика:*\n")
	builder.WriteString(fmt.Sprintf("💰 *Общая стоимость покупки:* $%f\n", totalBuyCost))
	builder.WriteString(fmt.Sprintf("🏦 *Текущая стоимость портфолио:* $%f\n", totalCurrentValue))
	builder.WriteString(fmt.Sprintf("%s *Общая прибыль/убыток:* $%f (%.2f%%)\n",
		totalProfitIcon, totalProfit, totalProfitPercent))

	b.sendMessage(update.Message.Chat.ID, builder.String())
}

func (b *Bot) HandleDeleteCrypto(update tgbotapi.Update) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		b.sendMessage(update.Message.Chat.ID, "Укажите ID записи: /delete_crypto <ID>\nID можно посмотреть в /my_wallet")
		return
	}

	id, err := strconv.Atoi(args)
	if err != nil {
		b.sendMessage(update.Message.Chat.ID, "Некорректный ID. Используйте числовой идентификатор.")
		return
	}

	err = b.walletService.DeleteCryptoFromWallet(id, update.Message.Chat.ID)
	if err != nil {
		log.Printf("Ошибка удаления из портфолио: %v", err)
		b.sendMessage(update.Message.Chat.ID, "Ошибка при удалении 😢")
		return
	}

	b.sendMessage(update.Message.Chat.ID, "✅ Запись успешно удалена из портфолио")
}
