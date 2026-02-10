"use client";

import { useState } from "react";
import { sendMail } from "@/lib/api";

interface ComposeFormProps {
  initialTo?: string;
  initialSubject?: string;
  initialBody?: string;
  threadId?: string;
  replyCode?: string;
  sendType?: "new" | "reply" | "forward";
  onClose: () => void;
  onSent: () => void;
}

export default function ComposeForm({
  initialTo = "",
  initialSubject = "",
  initialBody = "",
  threadId,
  replyCode,
  sendType = "new",
  onClose,
  onSent,
}: ComposeFormProps) {
  const [to, setTo] = useState(initialTo);
  const [subject, setSubject] = useState(initialSubject);
  const [body, setBody] = useState(initialBody);
  const [fromAddress, setFromAddress] = useState("");
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const recipients = to
      .split(",")
      .map((r) => r.trim())
      .filter(Boolean);

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
        send_type: sendType,
        thread_id: threadId,
        reply_code: replyCode,
        from_address: fromAddress || undefined,
      });
      onSent();
    } catch (err) {
      setError(err instanceof Error ? err.message : "送信に失敗しました");
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div className="bg-[var(--card-background)] rounded-xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col border border-[var(--card-border)]">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-[var(--card-border)]">
          <h2 className="font-semibold text-[var(--text-heading)]">
            {sendType === "new"
              ? "新規メール"
              : sendType === "reply"
              ? "返信"
              : "転送"}
          </h2>
          <button
            onClick={onClose}
            className="text-[var(--text-body)] hover:opacity-80 text-xl"
          >
            ×
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex-1 flex flex-col overflow-hidden">
          <div className="p-4 space-y-3 overflow-y-auto flex-1">
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
                placeholder="user@example.com"
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
            <div className="flex-1">
              <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
                本文
              </label>
              <textarea
                value={body}
                onChange={(e) => setBody(e.target.value)}
                rows={12}
                className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none resize-y"
                required
              />
            </div>
            {error && (
              <p className="text-sm text-red-600">{error}</p>
            )}
          </div>

          {/* Footer */}
          <div className="flex items-center justify-end gap-3 p-4 border-t border-[var(--card-border)]">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-[var(--text-body)] hover:opacity-80"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={sending}
              className="px-6 py-2 text-sm bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 disabled:opacity-50 transition-colors border border-[var(--card-border)]"
            >
              {sending ? "送信中..." : "送信"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
