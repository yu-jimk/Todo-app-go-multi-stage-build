package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"myapp/internal/db"

	"github.com/jackc/pgx/v5"
)

func main() {
	// 環境変数の取得とContextの設定
	dsn := os.Getenv("DB_SOURCE")
	if dsn == "" {
		log.Fatal("DB_SOURCE environment variable is required")
	}

	// 接続全体に10秒のタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// DB接続
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer conn.Close(ctx)

	// トランザクションの開始
	// 一つでも失敗したら「全部なかったこと」にするための仕組み
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("トランザクション開始失敗: %v", err)
	}
	defer tx.Rollback(ctx) // 正常終了（Commit）しなかった場合はロールバック

	// sqlcのインスタンス化（トランザクションを渡す）
	queries := db.New(conn).WithTx(tx)

	// 投入データ
	todos := []string{
		"牛乳を買う",
		"部屋の掃除をする",
		"Go言語の勉強をする",
	}

	log.Println("--- シードデータの投入を開始します ---")

	for _, title := range todos {
		inserted, err := queries.CreateTodo(ctx, db.CreateTodoParams{
			Title:     title,
			Completed: false,
		})
		if err != nil {
			// 一つでも失敗したらエラーとして終了（ロールバックされる）
			log.Fatalf("作成失敗 [%s]: %v", title, err)
		}
		fmt.Printf("✅ 作成成功: ID=%-3d Title=%s\n", inserted.ID, inserted.Title)
	}

	// 最後に全ての変更を確定させる
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("コミット失敗: %v", err)
	}

	log.Println("--- すべてのデータの投入が完了しました！ ---")
}