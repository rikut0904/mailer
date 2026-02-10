"use client";

import { useRouter } from "next/navigation";

export default function SettingsHomePage() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-[var(--background)] p-4">
      <div className="max-w-3xl mx-auto bg-[var(--card-background)] rounded-xl shadow p-6 border border-[var(--card-border)]">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl font-bold text-[var(--text-heading)]">設定</h1>
          <button
            onClick={() => router.push("/mail")}
            className="text-[var(--text-body)] hover:opacity-80"
          >
            ← 戻る
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <button
            onClick={() => router.push("/settings/notifications")}
            className="text-left p-4 rounded-lg border border-[var(--card-border)] bg-[var(--card-background)] hover:opacity-90 transition"
          >
            <p className="text-sm text-[var(--text-body)]">通知</p>
            <p className="text-base font-medium text-[var(--text-heading)]">
              Discord Webhook
            </p>
          </button>

          <button
            onClick={() => router.push("/settings/s3-domains")}
            className="text-left p-4 rounded-lg border border-[var(--card-border)] bg-[var(--card-background)] hover:opacity-90 transition"
          >
            <p className="text-sm text-[var(--text-body)]">受信</p>
            <p className="text-base font-medium text-[var(--text-heading)]">
              S3ドメイン設定
            </p>
          </button>
        </div>
      </div>
    </div>
  );
}
