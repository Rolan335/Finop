
## Api Финансовых операций

Стек: go, gin, pgx, postgresql, docker, goose.
Миграции выполняются внутри go.

## Запуск

Необходим Docker

```bash

  make run

```

Сервис будет локально доступен по адресу "localhost:8080"

Документация к Api доступна по адресу "localhost:8080/swagger"

Тесты запускаются командой "make test" ПОСЛЕ поднятых контейнеров

Пример .env файла находится в .env.example. Для Docker уже написан .env.docker