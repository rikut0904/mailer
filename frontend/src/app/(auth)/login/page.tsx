"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { signIn } = useAuth();
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await signIn(email, password);
      router.push("/mail");
    } catch {
      setError("メールアドレスまたはパスワードが正しくありません");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--background)]">
      <div className="bg-[var(--card-background)] p-8 rounded-xl shadow-lg w-full max-w-md border border-[var(--card-border)]">
        <h1 className="text-2xl font-bold text-[var(--text-heading)] mb-2">メール管理</h1>
        <p className="text-sm text-[var(--text-body)] mb-6">mailer</p>

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
          {error && <p className="text-sm text-red-600">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-[var(--primary-color)] text-[var(--text-heading)] rounded-lg hover:opacity-90 disabled:opacity-50 transition-colors font-medium border border-[var(--card-border)]"
          >
            {loading ? "サインイン中..." : "サインイン"}
          </button>
        </form>

        <div className="mt-4 text-sm text-[var(--text-body)]">
          初めての方はこちら
          <button
            onClick={() => router.push("/signup")}
            className="ml-2 text-[var(--primary-color)] hover:opacity-80"
          >
            新規登録
          </button>
        </div>
      </div>
    </div>
  );
}
