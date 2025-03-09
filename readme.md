# Собрать и запустить контейнеры
docker-compose up -d

# Остановить
docker-compose down

# Пересоздать БД (если нужно сбросить данные)
docker-compose down -v && docker-compose up -d