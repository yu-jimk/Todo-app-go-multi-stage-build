# ---------------------------------------------------
# 1. Base Stage: 共通の土台
# ---------------------------------------------------
FROM golang:1.25.3-alpine AS base

WORKDIR /app

# 証明書エラー対策（開発中も外部APIを呼ぶ可能性があるため）
RUN apk --no-cache add ca-certificates tzdata

# 日本時間(JST)に設定
ENV TZ=Asia/Tokyo

# go.mod / go.sum のみ先にコピーして依存解決
COPY go.mod go.sum ./
RUN go mod download

# ---------------------------------------------------
# 2. Dev Stage: 開発環境用
# ---------------------------------------------------
FROM base AS dev

# 開発用ツール (Gitなど)
RUN apk add --no-cache git make

# ホットリロードツール Air
RUN go install github.com/air-verse/air@v1.61.0

COPY . .

# .air.tomlが存在しない場合でも動くように、なければデフォルト設定で動かす
CMD ["air"]

# ---------------------------------------------------
# 3. Builder Stage: 本番ビルド用
# ---------------------------------------------------
FROM base AS builder

COPY . .

# 静的バイナリをビルド
# -trimpath を追加: ファイルパス情報を取り除き、セキュリティと再現性を向上
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o main ./cmd/api

# ---------------------------------------------------
# 4. Prod Stage: 本番実行環境
# ---------------------------------------------------
FROM alpine:3.20 AS prod

# 必要なパッケージ（証明書、タイムゾーン）をインストール
RUN apk --no-cache add ca-certificates tzdata

# Prod環境でもJSTにする
ENV TZ=Asia/Tokyo

# 非rootユーザーを作成
RUN adduser -D appuser

WORKDIR /app

# バイナリをコピー（所有権をappuserに変更しておく）
COPY --from=builder --chown=appuser:appuser /app/main .

USER appuser

CMD ["./main"]