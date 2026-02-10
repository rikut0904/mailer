"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { getUserSettings, updateUserSettings } from "@/lib/api";

export default function NotificationSettingsPage() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const [webhookUrl, setWebhookUrl] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
      return;
    }

    if (!authLoading && user) {
      (async () => {
        try {
          setLoading(true);
          const settings = await getUserSettings();
          setWebhookUrl(settings.discord_webhook_url || "");
        } catch (err) {
          setError(err instanceof Error ? err.message : "設定の取得に失敗しました");
        } finally {
          setLoading(false);
        }
      })();
    }
  }, [authLoading, user, router]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSaved(false);
    try {
      setSaving(true);
      const current = await getUserSettings();
      await updateUserSettings({
        discord_webhook_url: webhookUrl.trim(),
        selected_domain_id: current.selected_domain_id || "",
      });
      setSaved(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "保存に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  if (authLoading || !user || loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary-color)]" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[var(--background)] p-4">
      <div className="max-w-2xl mx-auto bg-[var(--card-background)] rounded-xl shadow p-6 border border-[var(--card-border)]">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl font-bold text-[var(--text-heading)]">通知設定</h1>
          <button
            onClick={() => router.push("/settings")}
            className="text-[var(--text-body)] hover:opacity-80"
          >
            ← 戻る
          </button>
        </div>

        <form onSubmit={handleSave} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              Discord Webhook URL
            </label>
            <input
              type="text"
              value={webhookUrl}
              onChange={(e) => setWebhookUrl(e.target.value)}
              placeholder="https://discord.com/api/webhooks/..."
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
            />
            <p className="mt-2 text-xs text-[var(--text-body)]">
              Discordのチャンネル設定 → 連携サービス → ウェブフックから作成できます。
            </p>
          </div>

          {error && <p className="text-sm text-red-600">{error}</p>}
          {saved && <p className="text-sm text-green-600">保存しました</p>}

          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => router.push("/settings")}
              className="px-4 py-2 text-[var(--text-body)] hover:opacity-80"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={saving}
              className="px-6 py-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 disabled:opacity-50 transition-colors font-medium border border-[var(--card-border)]"
            >
              {saving ? "保存中..." : "保存"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
