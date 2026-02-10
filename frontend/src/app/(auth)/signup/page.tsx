"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";

export default function SignupPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { signUp } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (password.length < 6) {
      setError("パスワードは6文字以上にしてください");
      return;
    }
    if (password !== confirm) {
      setError("パスワードが一致しません");
      return;
    }

    setLoading(true);
    try {
      await signUp(email, password);
      router.push("/mail");
    } catch {
      setError("アカウントの作成に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--background)]">
      <div className="bg-[var(--card-background)] p-8 rounded-xl shadow-lg w-full max-w-md border border-[var(--card-border)]">
        <h1 className="text-2xl font-bold text-[var(--text-heading)] mb-2">メール管理</h1>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              メールアドレス
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              パスワード
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-[var(--text-body)] mb-1">
              パスワード（確認）
            </label>
            <input
              type="password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              className="w-full px-3 py-2 border border-[var(--input-border)] bg-[var(--input-background)] rounded-lg focus:ring-2 focus:ring-[var(--primary-color)] focus:border-[var(--primary-color)] outline-none"
              required
            />
          </div>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 disabled:opacity-50 transition-colors font-medium border border-[var(--card-border)]"
          >
            {loading ? "作成中..." : "アカウント作成"}
          </button>
        </form>

        <div className="mt-4 text-sm text-[var(--text-body)]">
          すでにアカウントをお持ちですか？
          <button
            onClick={() => router.push("/login")}
            className="ml-2 text-[var(--primary-color)] hover:opacity-80"
          >
            サインイン
          </button>
        </div>
      </div>
    </div>
  );
}
