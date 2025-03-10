# my_anki_app
Подключиться с базе
docker exec -it spaced_repetition_db psql -U postgres -d spaced_repetition_db

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
Запуск контейнеров
make up
Остановка контейнеров
make down

После успешного запуска можно проверить доступность сервера командой:
curl http://localhost:8080/ping
Если сервер работает, вы получите ответ:
{"message": "pong"}



Структура проекта

my_app/
│── backend/              # Исходный код backend-сервиса
│   ├── cmd/              # Точка входа в приложение (main.go)
│   ├── internal/         # Внутренние пакеты
│   │   ├── db/           # Работа с базой данных
│   │   ├── handlers/     # HTTP-обработчики
│   │   ├── models/       # Описание структур данных
│   ├── migrations/       # Файлы миграций для БД
│   ├── pkg/              # Пакеты с вспомогательными функциями
│   ├── tests/            # Тесты для backend
│   ├── Dockerfile        # Docker-образ для backend
│   ├── go.mod            # Зависимости Go
│   ├── go.sum            # Контрольные суммы зависимостей
│── frontend/             # Исходный код frontend (React)
│── docker/               # Конфигурация контейнеров
│   ├── docker-compose.yml # Определение сервисов (PostgreSQL + backend)
│── docs/                 # Документация


Инициализация доки swagger. Из папки backend/cmd:
swag init -g main.go -o ../docs

