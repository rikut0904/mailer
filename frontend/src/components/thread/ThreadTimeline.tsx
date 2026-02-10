"use client";

import { useEffect, useState } from "react";
import { getThread } from "@/lib/api";
import type { ThreadResponse } from "@/types";

interface ThreadTimelineProps {
  threadId: string;
  onClose: () => void;
}

export default function ThreadTimeline({ threadId, onClose }: ThreadTimelineProps) {
  const [thread, setThread] = useState<ThreadResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchThread() {
      try {
        setLoading(true);
        const data = await getThread(threadId);
        setThread(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : "スレッドの取得に失敗しました");
      } finally {
        setLoading(false);
      }
    }
    fetchThread();
  }, [threadId]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary-color)]" />
      </div>
    );
  }

  if (error || !thread) {
    return (
      <div className="p-4">
        <button onClick={onClose} className="text-[var(--primary-color)] hover:underline mb-4">
          ← 戻る
        </button>
        <p className="text-red-600">{error || "スレッドが見つかりません"}</p>
      </div>
    );
  }

  return (
    <div className="p-4 md:p-6 max-w-4xl">
      <button onClick={onClose} className="text-[var(--primary-color)] hover:underline mb-4">
        ← 受信箱へ
      </button>

      <h2 className="text-xl font-bold text-[var(--text-heading)] mb-6">{thread.group_name}</h2>

      <div className="space-y-4">
        {thread.messages.map((msg, idx) => (
          <div
            key={idx}
            className={`
              rounded-lg p-4 border
              ${msg.type === "sent"
                ? "bg-[var(--primary-light)] border-[var(--card-border)] ml-8"
                : "bg-[var(--card-background)] border-[var(--card-border)] mr-8"
              }
            `}
          >
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <span
                  className={`text-xs px-2 py-0.5 rounded-full ${
                    msg.type === "sent"
                      ? "bg-[var(--primary-light)] text-[var(--text-heading)]"
                      : "bg-[var(--primary-light)] text-[var(--text-body)]"
                  }`}
                >
                  {msg.type === "sent" ? "送信" : "受信"}
                </span>
                <span className="text-sm font-medium text-[var(--text-body)]">
                  {msg.type === "sent" ? `宛先: ${msg.to}` : `差出人: ${msg.from}`}
                </span>
              </div>
              <span className="text-xs text-[var(--text-body)]">
                {new Date(msg.date).toLocaleString("ja-JP")}
              </span>
            </div>
            <p className="text-sm font-medium text-[var(--text-heading)] mb-1">{msg.subject}</p>
            <pre className="text-sm text-[var(--text-body)] whitespace-pre-wrap font-sans">
              {msg.body}
            </pre>
          </div>
        ))}
      </div>
    </div>
  );
}
