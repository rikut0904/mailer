"use client";

import { useState } from "react";
import type { ParsedMail } from "@/types";

interface MailDetailProps {
  mail: ParsedMail | null;
  onMarkRead: (s3Key: string, isRead: boolean) => void;
  onToggleStar: (s3Key: string, isStarred: boolean) => void;
  onDelete: (s3Key: string) => void;
  onReply: (mail: ParsedMail) => void;
  onForward: (mail: ParsedMail) => void;
  onViewThread: (threadId: string) => void;
}

export default function MailDetail({
  mail,
  onMarkRead,
  onToggleStar,
  onDelete,
  onReply,
  onForward,
  onViewThread,
}: MailDetailProps) {
  const [showHtml, setShowHtml] = useState(false);

  if (!mail) {
    return (
      <div className="flex items-center justify-center h-full text-[var(--text-body)]">
        メールを選択してください
      </div>
    );
  }

  return (
    <div className="p-4 md:p-6 max-w-4xl">
      {/* Header */}
      <div className="mb-4">
        <div className="flex items-start justify-between">
          <h1 className="text-xl font-bold text-[var(--text-heading)]">
            {mail.subject || "(件名なし)"}
          </h1>
          <button
            onClick={() => onToggleStar(mail.s3_key, !mail.is_starred)}
            className="text-2xl ml-2 flex-shrink-0"
          >
            {mail.is_starred ? "★" : "☆"}
          </button>
        </div>
        <div className="mt-2 text-sm text-[var(--text-body)]">
          <p>
            <span className="font-medium">差出人:</span> {mail.from}
          </p>
          <p>
            <span className="font-medium">宛先:</span> {mail.to}
          </p>
          <p>
            <span className="font-medium">日時:</span>{" "}
            {new Date(mail.date).toLocaleString("ja-JP")}
          </p>
          {mail.thread_id && (
            <button
              onClick={() => onViewThread(mail.thread_id!)}
              className="text-[var(--primary-color)] hover:underline mt-1"
            >
              スレッドを表示
            </button>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="flex flex-wrap gap-2 mb-4 pb-4 border-b border-[var(--card-border)]">
        <button
          onClick={() => onReply(mail)}
          className="px-3 py-1.5 text-sm bg-[var(--primary-color)] text-[var(--text-heading)] rounded hover:opacity-90 transition-colors border border-[var(--card-border)]"
        >
          返信
        </button>
        <button
          onClick={() => onForward(mail)}
          className="px-3 py-1.5 text-sm bg-[var(--primary-light)] text-[var(--text-heading)] rounded hover:opacity-90 transition-colors"
        >
          転送
        </button>
        <button
          onClick={() => onMarkRead(mail.s3_key, !mail.is_read)}
          className="px-3 py-1.5 text-sm bg-[var(--primary-light)] text-[var(--text-heading)] rounded hover:opacity-90 transition-colors"
        >
          {mail.is_read ? "未読にする" : "既読にする"}
        </button>
        <button
          onClick={() => onDelete(mail.s3_key)}
          className="px-3 py-1.5 text-sm bg-red-50 text-red-600 rounded hover:bg-red-100 transition-colors"
        >
          削除
        </button>
        {mail.html_body && (
          <button
            onClick={() => setShowHtml(!showHtml)}
            className="px-3 py-1.5 text-sm bg-[var(--primary-light)] text-[var(--text-heading)] rounded hover:opacity-90 transition-colors"
          >
            {showHtml ? "テキスト表示" : "HTML表示"}
          </button>
        )}
      </div>

      {/* Body */}
      {showHtml && mail.html_body ? (
        <div
          className="prose max-w-none"
          dangerouslySetInnerHTML={{ __html: mail.html_body }}
        />
      ) : (
        <pre className="whitespace-pre-wrap text-sm text-[var(--text-body)] font-sans leading-relaxed">
          {mail.body}
        </pre>
      )}

      {/* Attachments */}
      {mail.attachments && mail.attachments.length > 0 && (
        <div className="mt-6 pt-4 border-t border-[var(--card-border)]">
          <h3 className="text-sm font-medium text-[var(--text-body)] mb-2">
            添付ファイル ({mail.attachments.length})
          </h3>
          <div className="flex flex-wrap gap-2">
            {mail.attachments.map((att, idx) => (
              <div
                key={idx}
                className="px-3 py-2 bg-[var(--primary-light)] rounded text-sm text-[var(--text-body)]"
              >
                {att.filename} ({(att.size / 1024).toFixed(1)} KB)
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
