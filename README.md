# Расписание приёма лекарств

HTTP-сервис для хранения расписания приёма лекарств.
Go, SQLite.

## Запуск

go run .

## Пример запроса

curl -X POST http://localhost:8080/schedule \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"drug_name":"Аспирин","period":"2h","course_days":7}'
