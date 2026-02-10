# Mail Manager — rikut0904.site

独自ドメイン `rikut0904.site` のメール送受信を管理する個人用 Web アプリケーション。
Amazon SES / S3 を基盤に、受信メールの閲覧・返信・転送・スレッド管理を行う。

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| Backend | Go / Echo / GORM |
| Frontend | Next.js (App Router) / Tailwind CSS |
| Database | PostgreSQL |
| Auth | Firebase Authentication |
| Mail | Amazon SES (送信) / Amazon S3 (受信保存) |
| 通知 | Discord Webhook |

## クイックスタート

```bash
cp .env.example .env   # 値を埋めてから実行
docker compose up -d
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## ドキュメント

詳細な仕様は [`doc/`](./doc) を参照してください。

| ドキュメント | 内容 |
|------------|------|
| [project_overview.md](./doc/project_overview.md) | プロジェクト概要・インフラ前提 |
| [backend_spec.md](./doc/backend_spec.md) | バックエンド設計・API仕様 |
| [frontend_spec.md](./doc/frontend_spec.md) | フロントエンドUI/UX仕様 |
| [database_spec.md](./doc/database_spec.md) | データベーススキーマ |
| [thread_and_forward_spec.md](./doc/thread_and_forward_spec.md) | スレッド管理・転送ロジック |
