# Цитатник (QuoteAPI)

Мини-сервис для хранения и управления цитатами на Go.

---

## Функциональность

- Добавление новой цитаты (POST `/quotes`)
- Получение всех цитат (GET `/quotes`)
- Получение случайной цитаты (GET `/quotes/random`)
- Фильтрация цитат по автору (GET `/quotes?author=ИмяАвтора`)
- Удаление цитаты по ID (DELETE `/quotes/{id}`)

---

## Технологии

- Go
- Gorilla Mux
- Стандартная библиотека Go для работы с HTTP и JSON

---

## Запуск сервера

1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/yourusername/QuotesAPI.git
   cd QuotesAPI
   
2. Запустите сервер:
   ```bash
   go run cmd/main.go
3. Сервер будет доступен по адресу: http://localhost:8080
   
---

## Примеры запросов (curl)
  - Новая цитата : curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"YourQuote"}'

  - Получение всех цитат : curl http://localhost:8080/quotes
  - Получение случайной цитаты : curl http://localhost:8080/quotes/random
  - Получение цитат по автору : curl http://localhost:8080/quotes?author=Confucius
  - Удаление цитаты по ID : curl -X DELETE http://localhost:8080/quotes/1

