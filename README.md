# Расписание приёма лекарств

HTTP-сервис для хранения расписания приёма лекарств.
Go, SQLite, Docker, Docker Compose.

## Конфигурация

Настройки читаются из переменных окружения. Шаблон — `.env.example`:

В случае загрузки с Github отдельно нужно будет скопировать `.env.example` в `.env`  (compose подхватит его сам):

```powershell
Copy-Item .env.example .env
```

## Запуск

Локально:

```powershell
go run .
```

В Docker:

```powershell
docker compose up -d --build
```

- Данные живут в томе `db_data`: `docker compose down` их сохраняет, `docker compose down -v` — стирает.

## Пример запроса

Тело берём из файла `body.json` — одна и та же команда работает в PowerShell и cmd, кавычки ничего не ломают.

1. Создайте `body.json` в корне проекта:

```json
{"user_id":1,"drug_name":"Aspirin","period":"2h","course_days":7}
```

2. Отправьте:

```powershell
curl.exe -X POST "http://localhost:8080/schedule" -H "Content-Type: application/json" -d "@body.json"
```

3. Чтение (или просто URL в браузере):

```powershell
curl.exe "http://localhost:8080/schedule?user_id=1&schedule_id=1"
curl.exe "http://localhost:8080/next_takings?user_id=1"
```

## Образ в Docker Hub

Образ: `https://hub.docker.com/r/regreat39/medicine-schedule-app`

```powershell
docker pull regreat39/medicine-schedule-app:v1

docker run -d --name medicine -p 8080:8080 regreat39/medicine-schedule-app
```

## Просмотр БД

```powershell
docker compose exec app sh -c "apk add --no-cache sqlite && sqlite3 /app/data/medicine.db 'SELECT * FROM schedules;'"
```