package main

import (
	"fmt"
	"my_app_backend/internal/db"
)

func main() {
	fmt.Println("🚀 Запуск приложения...")
	db.ConnectDB()
}
