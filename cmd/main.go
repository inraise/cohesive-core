package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки файла .env")
	}
}