import { getAuth } from "./firebase";
import type {
  MailListResponse,
  ParsedMail,
  ThreadGroup,
  ThreadResponse,
  SendRequest,
  SendResponse,
  UserSettings,
  S3Domain,
} from "@/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function getAuthHeaders(): Promise<HeadersInit> {
  const auth = await getAuth();
  const user = auth.currentUser;
  if (!user) throw new Error("Not authenticated");

  const token = await user.getIdToken();
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const headers = await getAuthHeaders();
  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: { ...headers, ...options?.headers },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `API error: ${res.status}`);
  }

  return res.json();
}

export async function getMails(
  recipient?: string,
  page = 1,
  perPage = 20
): Promise<MailListResponse> {
  const params = new URLSearchParams();
  if (recipient) params.set("recipient", recipient);
  params.set("page", String(page));
  params.set("per_page", String(perPage));
  return apiFetch<MailListResponse>(`/api/mails?${params}`);
}

export async function getMail(s3Key: string): Promise<ParsedMail> {
  return apiFetch<ParsedMail>(`/api/mails/${encodeURIComponent(s3Key)}`);
}

export async function updateReadStatus(
  s3Key: string,
  isRead: boolean
): Promise<void> {
  await apiFetch(`/api/mails/${encodeURIComponent(s3Key)}/read`, {
    method: "PATCH",
    body: JSON.stringify({ is_read: isRead }),
  });
}

export async function updateStarStatus(
  s3Key: string,
  isStarred: boolean
): Promise<void> {
  await apiFetch(`/api/mails/${encodeURIComponent(s3Key)}/star`, {
    method: "PATCH",
    body: JSON.stringify({ is_starred: isStarred }),
  });
}

export async function deleteMail(s3Key: string): Promise<void> {
  await apiFetch(`/api/mails/${encodeURIComponent(s3Key)}`, {
    method: "DELETE",
  });
}

export async function syncMails(): Promise<{ synced: number }> {
  return apiFetch<{ synced: number }>("/api/mails/sync", { method: "POST" });
}

export async function getThreads(): Promise<ThreadGroup[]> {
  return apiFetch<ThreadGroup[]>("/api/threads");
}

export async function getThread(threadId: string): Promise<ThreadResponse> {
  return apiFetch<ThreadResponse>(`/api/threads/${encodeURIComponent(threadId)}`);
}

export async function sendMail(req: SendRequest): Promise<SendResponse> {
  return apiFetch<SendResponse>("/api/send", {
    method: "POST",
    body: JSON.stringify(req),
  });
}

export async function getUserSettings(): Promise<UserSettings> {
  return apiFetch<UserSettings>("/api/settings");
}

export async function updateUserSettings(
  settings: UserSettings
): Promise<UserSettings> {
  return apiFetch<UserSettings>("/api/settings", {
    method: "PUT",
    body: JSON.stringify(settings),
  });
}

export async function getRecipients(): Promise<{ recipients: string[] }> {
  return apiFetch<{ recipients: string[] }>("/api/mails/recipients");
}

export async function getDomains(): Promise<S3Domain[]> {
  return apiFetch<S3Domain[]>("/api/domains");
}

export async function createDomain(domain: Partial<S3Domain>): Promise<S3Domain> {
  return apiFetch<S3Domain>("/api/domains", {
    method: "POST",
    body: JSON.stringify(domain),
  });
}

export async function updateDomain(id: string, domain: Partial<S3Domain>): Promise<S3Domain> {
  return apiFetch<S3Domain>(`/api/domains/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(domain),
  });
}

export async function deleteDomain(id: string): Promise<void> {
  await apiFetch(`/api/domains/${encodeURIComponent(id)}`, { method: "DELETE" });
}
