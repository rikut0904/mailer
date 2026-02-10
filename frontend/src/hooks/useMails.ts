"use client";

import { useState, useEffect, useCallback } from "react";
import {
  getMails,
  updateReadStatus,
  updateStarStatus,
  deleteMail,
  syncMails,
} from "@/lib/api";
import type { ParsedMail, MailListResponse } from "@/types";

export function useMails(recipient?: string) {
  const [data, setData] = useState<MailListResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);

  const fetchMails = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await getMails(recipient, page);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch mails");
    } finally {
      setLoading(false);
    }
  }, [recipient, page]);

  useEffect(() => {
    fetchMails();
  }, [fetchMails]);

  const markRead = async (s3Key: string, isRead: boolean) => {
    await updateReadStatus(s3Key, isRead);
    setData((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        mails: prev.mails.map((m) =>
          m.s3_key === s3Key ? { ...m, is_read: isRead } : m
        ),
      };
    });
  };

  const toggleStar = async (s3Key: string, isStarred: boolean) => {
    await updateStarStatus(s3Key, isStarred);
    setData((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        mails: prev.mails.map((m) =>
          m.s3_key === s3Key ? { ...m, is_starred: isStarred } : m
        ),
      };
    });
  };

  const removeMail = async (s3Key: string) => {
    await deleteMail(s3Key);
    setData((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        mails: prev.mails.filter((m) => m.s3_key !== s3Key),
        total: prev.total - 1,
      };
    });
  };

  const sync = useCallback(async () => {
    const result = await syncMails();
    if (result.synced > 0) {
      await fetchMails();
    }
    return result;
  }, [fetchMails]);

  return {
    mails: data?.mails ?? [],
    total: data?.total ?? 0,
    page: data?.page ?? 1,
    totalPages: data?.total_pages ?? 1,
    loading,
    error,
    setPage,
    markRead,
    toggleStar,
    removeMail,
    sync,
    refresh: fetchMails,
  };
}
