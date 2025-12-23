# データベース・スキーマ (PostgreSQL)

## 1. mail_states (状態管理)
- `s3_key` (TEXT/PK): S3のオブジェクトキー
- `recipient_address` (TEXT): 受信先アドレス
- `is_read` (BOOLEAN): 既読フラグ（手動更新可）
- `is_starred` (BOOLEAN): スターフラグ
- `created_at` (TIMESTAMP)

## 2. thread_groups (スレッド親管理)
- `parent_uuid` (TEXT/PK): 親UUID
- `group_name` (TEXT): スレッド名（件名等）

## 3. sent_mails (送信履歴 & 1:N管理)
- `management_code` (TEXT/PK): 子UUID（本文挿入用）
- `parent_thread_id` (FK): thread_groupsへの参照
- `recipient_email` (TEXT): 送信先アドレス
