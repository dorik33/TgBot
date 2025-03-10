# Telegram Crypto Bot

Telegram бот для отслеживания цен криптовалют и управления подписками.

## Требования
- Docker и Docker Compose
- Go 1.16+
- API-токен от [BotFather](https://t.me/BotFather)

# Установка
1. Клонируйте репозиторий
2. Создайте `config.yaml` в корне проекта:

```yaml
botkey: "ВАШ_TELEGRAM_BOT_TOKEN"
timeout: 4s
database:
  host: localhost
  port: 5433
  user: bot_user
  password: secure_password
  dbname: crypto_bot
```
## Запуск postgres в Docker
docker-compose up -d

## Остановить
docker-compose down

## Пересоздать БД (сброс данных)
docker-compose down -v && docker-compose up -d

## Запуск бота
```go run cmd/main/main.go```

# Команды бота
```/price <токен>```	Текущая цена криптовалюты	```/price bitcoin```
```/sub <токен>```	Подписаться на обновления цены	```/sub ethereum```
```/subs```	Список активных подписок	
```/delete_sub <токен>```	Удалить подписку	```/delete_sub ethereum```
