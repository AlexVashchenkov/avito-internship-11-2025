# Тестовое задание Авито - сервис назначения ревьюеров для Pull Request-ов

## Инструкция по запуску сервиса:

1. Склонировать проект
```
git clone https://github.com/AlexVashchenkov/avito-internship-11-2025.git
cd avito-internship-11-2025
```

2. Создать `.env` файл по примеру из `.env.example`

3. Убедиться, что в файле `entry_point.sh` стоят LF-окончания строк (иначе Docker не подхватит этот файл и сервис не запустится)

4. Собрать и запустить сервис:
```
docker compose up --build
```

### Проблемы, с которыми столкнулся
1. При генерации кода с помощью утилиты ogen, столкнулся с проблемой, что в OpenAPI-спецификации не были указаны AdminToken и UserToken:
```
components:
    ...
    securitySchemes:
    AdminToken:
        type: http
        scheme: bearer
        bearerFormat: JWT
    UserToken:
        type: http
        scheme: bearer
        bearerFormat: JWT
```
2. При накатывании миграций выбрал утилиту `goose`, так как возникли большие проблемы с `golang-migrate`, для выполнения миграций через `docker-compose` установил её в сам образ и выполнил через скрипт `entry_point.sh`

### Выполненные доп. задания:
1. Описал конфигурацию линтера `golangci-lint` в файле `.golangci.yml`, и добавил команду для запуска линтера в `Makefile`

2. Имплементировал как in-memory реализацию хранилища, так и хранение данных в базе PostgreSQL

3. Добавил unit-тесты на in-memory хранилище, сгенерировал mock-структуры для последующих интеграционных тестов