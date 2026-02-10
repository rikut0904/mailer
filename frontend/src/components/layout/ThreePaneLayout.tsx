"use client";

import { ReactNode } from "react";

interface ThreePaneLayoutProps {
  sidebar: ReactNode;
  list: ReactNode;
  detail: ReactNode;
  showDetail: boolean;
  onBackToList?: () => void;
}

export default function ThreePaneLayout({
  sidebar,
  list,
  detail,
  showDetail,
  onBackToList,
}: ThreePaneLayoutProps) {
  return (
    <div className="flex h-screen bg-[var(--background)]">
      {/* Sidebar */}
      {sidebar}

      {/* List pane */}
      <div
        className={`
          w-full md:w-96 md:min-w-[24rem] border-r border-[var(--card-border)] bg-[var(--card-background)]
          overflow-y-auto flex-shrink-0
          ${showDetail ? "hidden md:block" : "block"}
        `}
      >
        {list}
      </div>

      {/* Detail pane */}
      <div
        className={`
          flex-1 bg-[var(--card-background)] overflow-y-auto
          ${showDetail ? "block" : "hidden md:block"}
        `}
      >
        {showDetail && onBackToList && (
          <button
            onClick={onBackToList}
            className="md:hidden p-3 text-[var(--primary-color)] hover:opacity-80 flex items-center gap-1"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            戻る
          </button>
        )}
        {detail}
      </div>
    </div>
  );
}
