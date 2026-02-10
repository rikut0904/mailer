# バックエンド設計・機能要件

## 1. クリーンアーキテクチャ構造
- `domain`: エンティティ定義、リポジトリ・インターフェース
- `usecase`: 業務シナリオ（メール解析フロー、UUID発行ロジック、削除フロー）
- `infrastructure`: AWS SDK v2, PostgreSQL, Discord連携の実装
- `interfaces`: HTTPハンドラー (Gin/Echo), Firebaseトークン検証

## 2. 主要機能
- **オンデマンド取得 & ページネーション**:
    - S3の `ListObjectsV2` を用い、一度に取得するキー数を制限（20件等）。
    - `ContinuationToken` を用いたページ遷移を実現。
- **メール解析 (MIME Parser)**:
    - S3から取得したRawデータを解析し、Subject/Body/From/Date/添付ファイルを抽出。
- **リソース管理 & 削除**:
    - 既読/未読/スター状態をDBで管理。ユーザーによる「未読への変更」を許可。
    - アプリ上の削除操作で、**DBレコードとS3オブジェクトを同時に物理削除**。
- **UUIDスレッド管理 (1:N対応)**:
    - 送信時に `【管理コード: UUID-XXXX】` を付与。
    - 1:N送信時は宛先ごとに個別の子UUIDを発行し、受信時に親UUIDへ紐付け。

## 3. 必要な環境変数
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` (ap-northeast-1)
- `S3_BUCKET_NAME`, `DATABASE_URL`
