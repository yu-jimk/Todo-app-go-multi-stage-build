# Todo-app-go-multi-stage-build

GoとMulti-stage-buildを使用したTodoバックエンド

## 環境セットアップ / シード / API確認

以下はローカル開発環境での手順です。`.env.example` をコピーして必要な値を設定してください。

1) `.env` を作る

```bash
cp .env.example .env
# 必要に応じてエディタで編集してください
```

2) コンテナ起動（バックグラウンド）

```bash
# 開発コンテナ
docker compose -f docker-compose.dev.yml up --build -d
# db のヘルスチェックが通るまで待ってください

# または　本番コンテナ
docker compose -f docker-compose.prod.yml up --build -d
```

3) スキーマをDBに適用

```bash
docker compose -f docker-compose.dev.yml exec -T db \
	psql -U ${DB_USER} -d ${DB_NAME} < sql/schema/schema.sql
```

4) シードデータ投入

```bash
# コンテナ内で実行（開発イメージに go が入っている想定）
docker compose -f docker-compose.dev.yml exec backend go run ./cmd/seed
```

5) API を curl で確認（サーバは http://localhost:8000 を想定）

```bash
# ヘルスチェック
curl http://localhost:8000/healthz
```

```bash
# 一覧取得
curl http://localhost:8000/todos
```

```bash
# 作成
curl -X POST http://localhost:8000/todos \
	-H "Content-Type: application/json" \
	-d '{"title":"買い物に行く"}'
```

```bash
# 単一取得
curl http://localhost:8000/todos/1
```

```bash
# タイトル更新
curl -X PATCH http://localhost:8000/todos/1/title \
	-H "Content-Type: application/json" \
	-d '{"title":"新しいタイトル"}'
```

```bash
# 完了状態更新
curl -X PATCH http://localhost:8000/todos/1/completed \
	-H "Content-Type: application/json" \
	-d '{"completed":true}'
```

```bash
# 削除
curl -X DELETE http://localhost:8000/todos/1
```
