"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { sendMail } from "@/lib/api";

export default function ComposePage() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const [to, setTo] = useState("");
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const [fromAddress, setFromAddress] = useState("");
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]);

  if (authLoading || !user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary-color)]" />
      </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    const recipients = to.split(",").map((r) => r.trim()).filter(Boolean);
    if (recipients.length === 0) {
      setError("宛先を1件以上入力してください");
      return;
    }

    if (!fromAddress.trim()) {
      setError("差出人を入力してください");
      return;
    }

    try {
      setSending(true);
      await sendMail({
        to: recipients,
        subject,
        body,
        send_type: "new",
        from_address: fromAddress || undefined,
      });
      router.push("/mail");
    } catch (err) {
      setError(err instanceof Error ? err.message : "送信に失敗しました");
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="min-h-screen bg-[var(--background)] p-4">
      <div className="max-w-2xl mx-auto bg-[var(--card-background)] rounded-xl shadow p-6 border border-[var(--card-border)]">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl font-bold text-[var(--text-heading)]">新規メール</h1>
          <button
            onClick={() => router.push("/mail")}
            className="text-[var(--text-body)] hover:opacity-80"
          >
            ← 戻る
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              差出人
            </label>
            <input
              type="email"
              value={fromAddress}
              onChange={(e) => setFromAddress(e.target.value)}
              placeholder="noreply@ml.rikut0904.site"
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              宛先（カンマ区切り）
            </label>
            <input
              type="text"
              value={to}
              onChange={(e) => setTo(e.target.value)}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              件名
            </label>
            <input
              type="text"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              本文
            </label>
            <textarea
              value={body}
              onChange={(e) => setBody(e.target.value)}
              rows={15}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none resize-y"
              required
            />
          </div>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => router.push("/mail")}
              className="px-4 py-2 text-[var(--text-body)] hover:opacity-80"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={sending}
              className="px-6 py-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 disabled:opacity-50 transition-colors font-medium border border-[var(--card-border)]"
            >
              {sending ? "送信中..." : "送信"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
