package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Docker Environment!")
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := "8000"
	fmt.Println("Server is running on port " + port)
	
    // DB接続情報の確認（接続確認はまだ実装していませんが、値が取れるか確認）
    dsn := os.Getenv("DB_SOURCE")
    fmt.Println("DB Connection String:", dsn)

	http.ListenAndServe(":"+port, nil)
}