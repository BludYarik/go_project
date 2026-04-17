package main

import (
	"context"
	"encoding/json" // Додали цей пакет для роботи з JSON
	"fmt"
	"net/http"
	"os"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	mode := os.Getenv("APP_MODE")
    fmt.Printf("Додаток запущено у режимі: %s\n", mode)
	// Отримуємо адресу Redis зі змінних оточення Docker
	redisAddr := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	http.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		// Збільшуємо лічильник у Redis
		visits, err := rdb.Incr(ctx, "visits").Result()
		
		if err != nil {
			fmt.Printf("ПОМИЛКА REDIS: %v\n", err)
			http.Error(w, "База даних недоступна", 500)
			return
		}

		// Створюємо структуру даних для відповіді
		response := map[string]interface{}{
			"visits": visits,
			"status": "success",
		}

		// Кажемо браузеру, що ми віддаємо саме JSON
		w.Header().Set("Content-Type", "application/json")
		
		// Перетворюємо нашу структуру в JSON-рядок і відправляємо
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Go-сервер запущено на порту 8080...")
	http.ListenAndServe(":8080", nil)
}
