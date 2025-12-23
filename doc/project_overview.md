# プロジェクト: 独自ドメインメール管理システム (mail.rikut0904.site)

## 1. プロジェクト概要
Amazon SES/S3を基盤とし、独自ドメイン `rikut0904.site` のメール送受信を管理する自分専用のWebアプリケーション。
「受信用アドレス」と「送信専用アドレス」を分離しつつ、独自の管理コード（UUID）を用いてスレッド管理を行う。
ユーザー操作によりDBメタデータだけでなく、AWS上の実体リソースも直接制御する。

## 2. 技術スタック
- **Backend**: Go (Railway) | アーキテクチャ: クリーンアーキテクチャ
- **Frontend**: Next.js (Vercel) | App Router, Tailwind CSS
- **Database**: PostgreSQL (Railway) | GORMによるORマッピング
- **Auth**: Firebase Authentication | 自分専用のアクセス制限
- **Infrastructure (AWS)**:
    - Amazon SES: 送信 (API利用)
    - Amazon S3: 受信メール (MIME形式) の保存・取得・削除
- **Integration**: Discord Webhook (着信通知)

## 3. インフラ前提
- SES 本番環境アクセス付与済み。
- カスタム MAIL FROM ドメイン (`mail.ml.rikut0904.site`) 設定済み。
- S3バケット (`rikut0904.site-ses`) にて受信メールを保存中。
