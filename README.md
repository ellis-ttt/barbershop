- лендинг: услуги и цены, мастера, галерея, контакты
- запись: услуга -> дата-> свободное время у мастера-> контакты-> подтверждение
- админка: записи на день и неделю, отмена, перенос, закрытие времени, статистика

Go (Gin, pgx, JWT), PostgreSQL, шаблоны Go + htmx + Alpine.js + Tailwind

## Запуск

```
docker compose up --build
```

- сайт: http://localhost:8080
- админка: http://localhost:8080/admin

## Настройки

Задаются переменными окружения в `docker-compose.yml`:

- `DATABASE_URL` - строка подключения к PostgreSQL
- `ADMIN_LOGIN`, `ADMIN_PASSWORD` - вход в админку
- `JWT_SECRET` - секрет для подписи токенов
- `TZ` - часовой пояс, по умолчанию `Europe/Moscow`

Перед деплоем замените `ADMIN_PASSWORD` и `JWT_SECRET`.
