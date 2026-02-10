"use client";

import type { ParsedMail } from "@/types";

interface MailListProps {
  mails: ParsedMail[];
  selectedKey: string | null;
  onSelect: (mail: ParsedMail) => void;
  onToggleStar: (s3Key: string, isStarred: boolean) => void;
  loading: boolean;
  page: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  onSync: () => void;
}

export default function MailList({
  mails,
  selectedKey,
  onSelect,
  onToggleStar,
  loading,
  page,
  totalPages,
  onPageChange,
  onSync,
}: MailListProps) {
  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary-color)]" />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-3 border-b border-[var(--card-border)] flex items-center justify-between">
        <h2 className="font-semibold text-[var(--text-heading)]">受信箱</h2>
        <button
          onClick={onSync}
          className="text-sm px-3 py-1 bg-[var(--primary-light)] text-[var(--text-heading)] hover:opacity-90 rounded transition-colors"
          title="同期"
        >
          ↻ 同期
        </button>
      </div>

      {/* Mail items */}
      <div className="flex-1 overflow-y-auto">
        {mails.length === 0 ? (
          <p className="p-4 text-center text-[var(--text-body)]">メールがありません</p>
        ) : (
          mails.map((mail) => (
            <div
              key={mail.s3_key}
              onClick={() => onSelect(mail)}
              className={`
                p-3 border-b border-[var(--card-border)] cursor-pointer transition-colors
                ${selectedKey === mail.s3_key ? "bg-[var(--primary-light)]" : "hover:bg-[var(--primary-light)]/50"}
                ${!mail.is_read ? "bg-[var(--card-background)]" : "bg-[var(--primary-light)]/30"}
              `}
            >
              <div className="flex items-start gap-2">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onToggleStar(mail.s3_key, !mail.is_starred);
                  }}
                  className="mt-1 text-lg flex-shrink-0"
                >
                  {mail.is_starred ? "★" : "☆"}
                </button>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between">
                    <p
                      className={`text-sm truncate ${
                        !mail.is_read ? "font-bold text-[var(--text-heading)]" : "text-[var(--text-body)]"
                      }`}
                    >
                      {mail.from}
                    </p>
                    <span className="text-xs text-[var(--text-body)] flex-shrink-0 ml-2">
                      {new Date(mail.date).toLocaleDateString("ja-JP", {
                        month: "short",
                        day: "numeric",
                      })}
                    </span>
                  </div>
                  <p
                    className={`text-sm truncate ${
                      !mail.is_read ? "font-semibold text-[var(--text-heading)]" : "text-[var(--text-body)]"
                    }`}
                  >
                    {mail.subject || "(件名なし)"}
                  </p>
                  <p className="text-xs text-[var(--text-body)] truncate mt-0.5">
                    {mail.body?.substring(0, 80)}
                  </p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="p-3 border-t border-[var(--card-border)] flex items-center justify-between">
          <button
            onClick={() => onPageChange(page - 1)}
            disabled={page <= 1}
            className="px-3 py-1 text-sm bg-[var(--primary-light)] text-[var(--text-heading)] hover:opacity-90 rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            前へ
          </button>
          <span className="text-sm text-[var(--text-body)]">
            {page} / {totalPages}
          </span>
          <button
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
            className="px-3 py-1 text-sm bg-[var(--primary-light)] text-[var(--text-heading)] hover:opacity-90 rounded disabled:opacity-50 disabled:cursor-not-allowed"
          >
            次へ
          </button>
        </div>
      )}
    </div>
  );
}
