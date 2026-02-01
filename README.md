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

## テスト

### ユニットテスト

ユニットテストは通常の `go test` で実行できます:

```bash
go test ./...
```

追加したテスト:
- `internal/service` : サービス層のユニットテスト


## CI / CD

このリポジトリには GitHub Actions のワークフローを追加しています。

- CI: [.github/workflows/ci.yml](.github/workflows/ci.yml#L1)
	- トリガー: `push`（`main` ブランチ）および `pull_request`
	- 内容: `actions/setup-go` を使い、`go mod download` → `gofmt` チェック → `go vet` → `go test ./... -v -race` を実行します。
	- 目的: コード整形・静的解析・ユニットテストの自動実行

- CD: [.github/workflows/cd.yml](.github/workflows/cd.yml#L1)
	- トリガー: `push`（`main` ブランチ）
	- 内容: マルチアーキテクチャ対応でイメージをビルドし、GitHub Container Registry (GHCR) にプッシュします。イメージメタ情報は `docker/metadata-action` で生成され、`docker/build-push-action` で `ghcr.io/${{ github.repository }}` にタグ付きでプッシュします。
	- 必要権限: ワークフロー内で `packages: write` を要求しており、`GITHUB_TOKEN` を用いて GHCR にログインします。
	- 追加デプロイ: オプションで SSH を使ったリモートデプロイが行えます（`appleboy/scp-action` と `appleboy/ssh-action` を使用）。

必要なシークレット（リポジトリの Settings → Secrets に追加）:

- `SSH_HOST` — デプロイ先ホスト（SSH）
- `SSH_USER` — デプロイ用ユーザー名
- `SSH_PRIVATE_KEY` — 秘密鍵（PEM 形式）
- `SSH_PORT` — （省略可、デフォルト 22）

メモ:

- デフォルトの CD は `ghcr.io/<owner>/<repo>:latest` などのタグでイメージをプッシュします。タグ方針を変えたい場合は `.github/workflows/cd.yml` の `docker/metadata-action` 設定を編集してください。
- SSH デプロイを有効にするには上記シークレットを設定してください。ワークフローはリモートの `/var/www/myapp` に `docker-compose.prod.yml` をコピーし、`docker compose -f docker-compose.prod.yml pull` と `docker compose -f docker-compose.prod.yml up -d` を実行します。必要に応じてパスやコマンドを編集してください。
- ワークフロー定義ファイル:
	- [.github/workflows/ci.yml](.github/workflows/ci.yml#L1)
	- [.github/workflows/cd.yml](.github/workflows/cd.yml#L1)
	- [docker-compose.prod.yml](docker-compose.prod.yml#L1)

まずは PR を作って CI を確認してください。問題なければ `main` にマージして CD を試すことができます。


