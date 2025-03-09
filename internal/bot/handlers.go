package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет я бот для отслеживания криптовалюты напиши /help чтобы увидеть список поддерживаемых комманд")
	bot.Send(msg)
}

func HelpMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(`Я могу помочь вам с основными вопросами!
Напишите /start для начала.
/price <название криптовалюты> - Отправляет текущую цену указанного токена.
/sub <название криптовалюты> - Позволяет добавить подписку на отслеживание цены криптовалюты.
/subs - Показывает все текущие подписки.
/delete_sub <название криптовалюты> - Позволяет удалить подписку на конкретный токен.`))
	bot.Send(msg)
}

func PriceMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, apiClient *api.APIClient) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	args = strings.ToLower(args)
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите название криптовалюты: /price <имя>")
		bot.Send(msg)
		return
	}

	crypto, err := apiClient.GetInfo(args)
	if err != nil {
		log.Printf("Ошибка при запросе цены: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные 😢")
		bot.Send(msg)
		return
	}

	price := fmt.Sprintf("💰 Цена %s (%s): $%f", crypto.ID, crypto.Symbol, parsePrice(crypto.PriceUSD))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, price)
	bot.Send(msg)
}

func AddSubscriptionMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, apiClient *api.APIClient, repo *database.SubscriptionRepository) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	args = strings.ToLower(args)
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите название криптовалюты: /sub <имя>")
		bot.Send(msg)
		return
	}
	crypto, err := apiClient.GetInfo(args)
	if err != nil {
		log.Printf("Ошибка при добавлении подписки: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные 😢")
		bot.Send(msg)
		return
	}
	price := fmt.Sprintf("💰 Цена %s (%s): $%f", crypto.ID, crypto.Symbol, parsePrice(crypto.PriceUSD))
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, price)
	bot.Send(msg)

	if err = repo.AddSubscription(update.Message.Chat.ID, args); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось добавить подписку 😢")
		bot.Send(msg)
		return
	}
	log.Printf("Подписка добавлена user: %d, token: %s\n", update.Message.Chat.ID, args)
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Подписка на %s успешно добавлена!", args))
	bot.Send(msg)
}

func GetSubcriptionsMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, apiClient *api.APIClient, repo *database.SubscriptionRepository) {
	subs, err := repo.GetSubcriptions(update.Message.Chat.ID)
	if err != nil {
		log.Printf("Ошибка при нахождении подписок: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить данные о подписках😢")
		bot.Send(msg)
	}
	log.Printf("Подписка получены user: %d, \n", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, parseSubs(subs))
	bot.Send(msg)
}

func DeleteSubcriptionsMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, apiClient *api.APIClient, repo *database.SubscriptionRepository) {
	args := strings.TrimSpace(update.Message.CommandArguments())
	args = strings.ToLower(args)
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите название криптовалюты: /price <имя>")
		bot.Send(msg)
		return
	}

	err := repo.DeleteSubscription(update.Message.Chat.ID, args)
	if err != nil {
		log.Println("Ошибка при удалении подписки: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось удалить подписку😢")
		bot.Send(msg)
	}
	log.Printf("Подписка удалены user: %d, \n", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Подписка на %s успешно удалена!", args))
	bot.Send(msg)
}
