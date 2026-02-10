"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { createDomain, deleteDomain, getDomains } from "@/lib/api";
import type { S3Domain } from "@/types";

export default function S3DomainsPage() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const [domains, setDomains] = useState<S3Domain[]>([]);
  const [domainForm, setDomainForm] = useState({
    name: "",
    bucket: "",
    region: "",
    access_key_id: "",
    secret_key: "",
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
      return;
    }

    if (!authLoading && user) {
      (async () => {
        try {
          setLoading(true);
          const domainList = await getDomains();
          setDomains(domainList);
        } finally {
          setLoading(false);
        }
      })();
    }
  }, [authLoading, user, router]);

  const handleCreateDomain = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      const created = await createDomain(domainForm);
      setDomains((prev) => [...prev, created]);
      setDomainForm({
        name: "",
        bucket: "",
        region: "",
        access_key_id: "",
        secret_key: "",
      });
      setMessage("ドメインを追加しました");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "追加に失敗しました");
    }
  };

  const handleDeleteDomain = async (id: string) => {
    setMessage(null);
    try {
      await deleteDomain(id);
      setDomains((prev) => prev.filter((d) => d.id !== id));
      setMessage("ドメインを削除しました");
    } catch (err) {
      setMessage(err instanceof Error ? err.message : "削除に失敗しました");
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
      <div className="max-w-3xl mx-auto bg-[var(--card-background)] rounded-xl shadow p-6 border border-[var(--card-border)]">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl font-bold text-[var(--text-heading)]">S3ドメイン設定</h1>
          <button
            onClick={() => router.push("/settings")}
            className="text-[var(--text-body)] hover:opacity-80"
          >
            ← 戻る
          </button>
        </div>

        <form onSubmit={handleCreateDomain} className="space-y-3 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            <input
              type="text"
              value={domainForm.name}
              onChange={(e) => setDomainForm({ ...domainForm, name: e.target.value })}
              placeholder="表示名"
              className="px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg"
              required
            />
            <input
              type="text"
              value={domainForm.bucket}
              onChange={(e) => setDomainForm({ ...domainForm, bucket: e.target.value })}
              placeholder="バケット名"
              className="px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg"
              required
            />
            <input
              type="text"
              value={domainForm.region}
              onChange={(e) => setDomainForm({ ...domainForm, region: e.target.value })}
              placeholder="リージョン"
              className="px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg"
              required
            />
            <input
              type="text"
              value={domainForm.access_key_id}
              onChange={(e) => setDomainForm({ ...domainForm, access_key_id: e.target.value })}
              placeholder="Access Key ID"
              className="px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg"
              required
            />
            <input
              type="password"
              value={domainForm.secret_key}
              onChange={(e) => setDomainForm({ ...domainForm, secret_key: e.target.value })}
              placeholder="Secret Access Key"
              className="px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg"
              required
            />
          </div>
          <button
            type="submit"
            className="px-4 py-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 border border-[var(--card-border)]"
          >
            追加
          </button>
        </form>

        {message && <p className="text-sm text-[var(--text-body)] mb-3">{message}</p>}

        <div className="space-y-2">
          {domains.map((d) => (
            <div
              key={d.id}
              className="flex items-center justify-between px-3 py-2 border border-[var(--card-border)] rounded-lg"
            >
              <div>
                <p className="text-sm font-medium text-[var(--text-heading)]">{d.name}</p>
                <p className="text-xs text-[var(--text-body)]">
                  {d.bucket} / {d.region}
                </p>
              </div>
              <button
                onClick={() => handleDeleteDomain(d.id)}
                className="text-sm text-red-600 hover:opacity-80"
              >
                削除
              </button>
            </div>
          ))}
          {domains.length === 0 && (
            <p className="text-sm text-[var(--text-body)]">ドメインが未登録です</p>
          )}
        </div>
      </div>
    </div>
  );
}
