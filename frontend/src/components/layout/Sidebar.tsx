"use client";

import { useState } from "react";

interface SidebarProps {
  currentRecipient: string;
  recipients: string[];
  domains: { id: string; name: string }[];
  selectedDomainId: string;
  onRecipientChange: (recipient: string) => void;
  onDomainChange: (domainId: string) => void;
  onCompose: () => void;
  onSettings: () => void;
  onSignOut: () => void;
}

export default function Sidebar({
  currentRecipient,
  recipients,
  domains,
  selectedDomainId,
  onRecipientChange,
  onDomainChange,
  onCompose,
  onSettings,
  onSignOut,
}: SidebarProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      {/* Mobile hamburger */}
      <button
        className="md:hidden fixed top-4 left-4 z-50 p-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg shadow border border-[var(--card-border)]"
        onClick={() => setIsOpen(!isOpen)}
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {isOpen ? (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          ) : (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          )}
        </svg>
      </button>

      {/* Overlay for mobile */}
      {isOpen && (
        <div
          className="md:hidden fixed inset-0 bg-black/50 z-30"
          onClick={() => setIsOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`
          fixed md:static inset-y-0 left-0 z-40
          w-64 bg-[var(--primary-color)] text-[var(--text-heading)] flex flex-col border-r border-[var(--card-border)]
          transform transition-transform duration-200
          ${isOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"}
        `}
      >
        <div className="p-4 border-b border-[var(--card-border)]">
          <h1 className="text-lg font-bold">メール管理</h1>
          <p className="text-xs text-[var(--text-body)]">rikut0904.site</p>
        </div>

        <div className="p-4">
          <button
            onClick={() => {
              onCompose();
              setIsOpen(false);
            }}
            className="w-full py-2 px-4 bg-[var(--card-background)] text-[var(--text-heading)] hover:opacity-90 rounded-lg font-medium transition-colors border border-[var(--card-border)]"
          >
            + 新規メール
          </button>
        </div>

        <nav className="flex-1 px-2">
          <p className="px-2 py-1 text-xs text-[var(--text-body)] uppercase tracking-wider">
            ドメイン
          </p>
          <div className="px-2 mb-3">
            <select
              value={selectedDomainId}
              onChange={(e) => {
                onDomainChange(e.target.value);
                setIsOpen(false);
              }}
              className="w-full px-3 py-2 text-sm rounded-lg border border-[var(--card-border)] bg-[var(--card-background)] text-[var(--text-body)]"
            >
              {domains.map((domain) => (
                <option key={domain.id} value={domain.id}>
                  {domain.name}
                </option>
              ))}
            </select>
          </div>

          <p className="px-2 py-1 text-xs text-[var(--text-body)] uppercase tracking-wider">
            受信アドレス
          </p>
          <button
            onClick={() => {
              onRecipientChange("");
              setIsOpen(false);
            }}
            className={`
              w-full text-left px-3 py-2 rounded-lg mb-1 transition-colors
              ${currentRecipient === ""
                ? "bg-[var(--card-background)] text-[var(--text-heading)]"
                : "text-[var(--text-body)] hover:bg-[var(--card-background)]"
              }
            `}
          >
            すべて
          </button>
          {recipients.map((addr) => (
            <button
              key={addr}
              onClick={() => {
                onRecipientChange(addr);
                setIsOpen(false);
              }}
              className={`
                w-full text-left px-3 py-2 rounded-lg mb-1 transition-colors
                ${currentRecipient === addr
                  ? "bg-[var(--card-background)] text-[var(--text-heading)]"
                  : "text-[var(--text-body)] hover:bg-[var(--card-background)]"
                }
              `}
            >
              {addr}
            </button>
          ))}
        </nav>

        <div className="p-4 border-t border-[var(--card-border)]">
          <button
            onClick={() => {
              onSettings();
              setIsOpen(false);
            }}
            className="w-full py-2 px-4 text-sm text-[var(--text-body)] hover:opacity-80 transition-colors"
          >
            設定
          </button>
          <button
            onClick={onSignOut}
            className="w-full py-2 px-4 text-sm text-[var(--text-body)] hover:opacity-80 transition-colors"
          >
            サインアウト
          </button>
        </div>
      </aside>
    </>
  );
}
