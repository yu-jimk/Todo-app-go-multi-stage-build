package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"myapp/internal/db"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/service"
)

func main() {
	// 設定の読み込み
	port := "8000"
	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		log.Fatal("DB_SOURCE is required")
	}

	// DB接続 (pgxpool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	// DI (依存性の注入) の組み立て
	queries := db.New(pool)
	todoRepo := repository.NewTodoRepository(queries)
	todoSvc := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoSvc)

	// ルーティング設定
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, todoHandler)

	// その他のルート
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Docker Environment!")
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// サーバー起動
	fmt.Println("Server is running on port " + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}