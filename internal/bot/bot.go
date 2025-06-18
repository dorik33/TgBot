package bot

import (
	"strings"

	"github.com/dorik33/TgBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	botAPI              *tgbotapi.BotAPI
	subscriptionService service.SubscriptionService
	walletService       service.WalletService
	cryptoService       service.CryptoService
}


func NewBot(botAPI *tgbotapi.BotAPI, subService service.SubscriptionService, walletService service.WalletService, cryptoService service.CryptoService) *Bot {
	return &Bot{
		botAPI:              botAPI,
		subscriptionService: subService,
		walletService:       walletService,
		cryptoService:       cryptoService,
	}
}

func (b *Bot) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}
	text := update.Message.Text
	command := strings.Split(text, " ")[0]

	switch command {
	case "/start":
		b.sendMessage(update.Message.Chat.ID, "Привет, я бот для отслеживания криптовалюты, напиши /help чтобы увидеть список поддерживаемых команд")
	case "/help":
		b.sendMessage(update.Message.Chat.ID, `*Доступные команды:*
	/start - Начало работы
	/price <токен> - Текущая цена криптовалюты
	/add_crypto <токен> <количество> [цена] - Добавить в портфолио
	/my_wallet - Показать мое портфолио
	/delete_crypto <ID> - Удалить из портфолио
	/sub <токен> - Подписаться на обновления
	/subs - Мои подписки
	/delete_sub <токен> - Отписаться`)
	case "/price":
		b.handlePrice(update)
	case "/sub":
		b.handleSub(update)
	case "/subs":
		b.handleListSubs(update)
	case "/delete_sub":
		b.handleDeleteSub(update)
	case "/add_crypto":
		b.HandleAddCrypto(update)
	case "/my_wallet":
		b.HandleMyWallet(update)
	case "/delete_crypto":
		b.HandleDeleteCrypto(update)
	default:
		b.sendMessage(update.Message.Chat.ID, "Неизвестная команда. Напиши /help")
	}
}
