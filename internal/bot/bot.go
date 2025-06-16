package bot

import (
	"log"
	"strings"

	"github.com/dorik33/TgBot/internal/api"
	"github.com/dorik33/TgBot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	apiClient  *api.APIClient
	supbrepo   *database.SubscriptionRepository
	walletrepo *database.WalletRepository
	botAPI     *tgbotapi.BotAPI
}

func NewBot(token string, apiClient *api.APIClient, subrepo *database.SubscriptionRepository, walletrepo *database.WalletRepository) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		apiClient:  apiClient,
		supbrepo:   subrepo,
		walletrepo: walletrepo,
		botAPI:     botAPI,
	}, nil
}

func (b *Bot) Start() {
	log.Println("Bot is running...")

	go b.ticker()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s, %d\n", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	text := update.Message.Text
	command := strings.Split(text, " ")[0]

	switch command {
	case "/start":
		b.sendMessage(update.Message.Chat.ID, "Привет я бот для отслеживания криптовалюты, напиши /help чтобы увидеть список поддерживаемых команд")
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
