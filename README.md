Бот для голосований в Mattermost с хранением данных в Tarantool.

## Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/mattermost-poll-bot.git
   ```

2. Запустите сервисы:
   ```bash
   docker-compose up --build
   ```

## Команды
- `!poll create <id> <option1> <option2>` — Создать голосование.
- `!poll vote <id> <option>` — Проголосовать.
- `!poll results <id>` — Показать результаты.

## Логирование
Логи пишутся в `stdout` в формате JSON. Уровни:
- `INFO` — Основные события.
- `ERROR` — Критические ошибки.