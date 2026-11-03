# Расписание приёма лекарств

HTTP-сервис для хранения расписания приёма лекарств.  
Go, SQLite, Docker, Docker-Compose.

## Запуск

Локально:
```powershell
go run .
```
В Docker:
```powershell
docker compose up
```
## Пример запроса
```powershell
curl -X POST http://localhost:8080/schedule -H "Content-Type: application/json" -d '{"user_id":1,"drug_name":"Aspirin","period":"2h","course_days":7}'
```
## Просмотр БД
Прямо в контейнере:
```powershell
docker compose exec app sh -c "apk add --no-cache sqlite && sqlite3 /app/data/medicine.db 'SELECT * FROM schedules;'"
```