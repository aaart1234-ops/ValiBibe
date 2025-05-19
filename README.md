# my_anki_app
Подключиться с базе
docker exec -it spaced_repetition_db psql -U postgres -d spaced_repetition_db

или к системной базе так
docker exec -it spaced_repetition_db psql -U postgres -d postgres
Посмотри список баз данных:
\l

Перезапуск контейнеров  
docker-compose down && docker-compose up -d --build

Запуск контейнеров с указанием пути к .env
docker-compose --env-file ../.env up -d --build


Найти процесс, занимающий порт 8080, и завершить его
lsof -i :8080
Заверши процесс:
kill -9 12345
Теперь можно заново запустить сервер:
go run cmd/main.go

Запуск логов backend:
make logs
Применение миграций:
make migrate
Создать миграцию:
make create-migration NAME=create_notes_table
Запуск контейнеров:
make up
Остановка контейнеров:
make down

После успешного запуска можно проверить доступность сервера командой:
curl http://localhost:8080/ping
Если сервер работает, вы получите ответ:
{"message": "pong"}

Инициализация доки swagger. Из папки backend:
swag init -g cmd/main.go -o docs

Запуск тестов из папки backend:
go test -v ./tests/...


