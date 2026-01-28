## Todo-app (Docker)

このリポジトリはGo製TodoアプリのマルチステージDockerビルド設定を含みます。

開発用と本番用それぞれの`docker-compose`設定が用意されています。

開発 (ホットリロード):

```bash
# ビルドしてバックエンドとDBを起動
docker compose -f docker-compose.dev.yml up --build

# 停止
docker compose -f docker-compose.dev.yml down
```

本番（ローカル検証）:

```bash
# 本番用イメージをビルドして起動（ポート80）
docker compose -f docker-compose.prod.yml up --build -d

# ログ確認
docker compose -f docker-compose.prod.yml logs -f

# 停止
docker compose -f docker-compose.prod.yml down
```

便利な`Makefile`を用意しています（`make help` を参照）。

環境変数は `.env.example` を参考に設定してください。

# Todo-app-go-multi-stage-build

GoとMulti-stage-buildを使用したTodoバックエンド
