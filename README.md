# Тестовое задание Авито - сервис назначения ревьюеров для Pull Request-ов

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