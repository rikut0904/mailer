"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { useMails } from "@/hooks/useMails";
import ThreePaneLayout from "@/components/layout/ThreePaneLayout";
import Sidebar from "@/components/layout/Sidebar";
import MailList from "@/components/mail/MailList";
import MailDetail from "@/components/mail/MailDetail";
import ComposeForm from "@/components/compose/ComposeForm";
import ThreadTimeline from "@/components/thread/ThreadTimeline";
import type { ParsedMail } from "@/types";
import { getDomains, getRecipients, getUserSettings, updateUserSettings } from "@/lib/api";

export default function MailPage() {
  const { user, loading: authLoading, signOut } = useAuth();
  const router = useRouter();
  const [recipient, setRecipient] = useState("");
  const [selectedMail, setSelectedMail] = useState<ParsedMail | null>(null);
  const [showCompose, setShowCompose] = useState(false);
  const [viewingThread, setViewingThread] = useState<string | null>(null);
  const [domains, setDomains] = useState<{ id: string; name: string }[]>([]);
  const [selectedDomainId, setSelectedDomainId] = useState("");
  const [recipients, setRecipients] = useState<string[]>([]);

  // Compose state for reply/forward
  const [composeProps, setComposeProps] = useState<{
    initialTo?: string;
    initialSubject?: string;
    initialBody?: string;
    threadId?: string;
    replyCode?: string;
    sendType: "new" | "reply" | "forward";
  }>({ sendType: "new" });

  const {
    mails,
    page,
    totalPages,
    loading,
    setPage,
    markRead,
    toggleStar,
    removeMail,
    sync,
    refresh,
  } = useMails(recipient);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]);

  useEffect(() => {
    if (!authLoading && user) {
      (async () => {
        try {
          const [domainList, settings] = await Promise.all([
            getDomains(),
            getUserSettings(),
          ]);
          setDomains(domainList.map((d) => ({ id: d.id, name: d.name })));
          const domainId = settings.selected_domain_id || domainList[0]?.id || "";
          setSelectedDomainId(domainId);
          if (domainId && domainId !== settings.selected_domain_id) {
            await updateUserSettings({
              discord_webhook_url: settings.discord_webhook_url || "",
              selected_domain_id: domainId,
            });
          }
          const recipientRes = await getRecipients();
          setRecipients(recipientRes.recipients || []);
          if (domainId) {
            await sync();
            await refresh();
          }
        } catch (err) {
          console.error(err);
        }
      })();
    }
  }, [authLoading, user, sync, refresh]);

  if (authLoading || !user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[var(--primary-color)]" />
      </div>
    );
  }

  const handleSelectMail = (mail: ParsedMail) => {
    setSelectedMail(mail);
    setViewingThread(null);
    if (!mail.is_read) {
      markRead(mail.s3_key, true);
    }
  };

  const handleReply = (mail: ParsedMail) => {
    setComposeProps({
      initialTo: mail.from,
      initialSubject: `Re: ${mail.subject}`,
      initialBody: `\n\n---\n${new Date(mail.date).toLocaleString("ja-JP")} に ${mail.from} さんが書きました:\n${mail.body}`,
      threadId: mail.thread_id || undefined,
      sendType: "reply",
    });
    setShowCompose(true);
  };

  const handleForward = (mail: ParsedMail) => {
    setComposeProps({
      initialSubject: `Fwd: ${mail.subject}`,
      initialBody: `\n\n---\n転送メッセージ:\n差出人: ${mail.from}\n日時: ${new Date(mail.date).toLocaleString("ja-JP")}\n件名: ${mail.subject}\n\n${mail.body}`,
      threadId: mail.thread_id || undefined,
      sendType: "forward",
    });
    setShowCompose(true);
  };

  const handleCompose = () => {
    setComposeProps({ sendType: "new" });
    setShowCompose(true);
  };

  const handleDelete = async (s3Key: string) => {
    if (!confirm("このメールを削除しますか？")) return;
    await removeMail(s3Key);
    if (selectedMail?.s3_key === s3Key) {
      setSelectedMail(null);
    }
  };

  const detailPane = viewingThread ? (
    <ThreadTimeline
      threadId={viewingThread}
      onClose={() => setViewingThread(null)}
    />
  ) : (
    <MailDetail
      mail={selectedMail}
      onMarkRead={markRead}
      onToggleStar={toggleStar}
      onDelete={handleDelete}
      onReply={handleReply}
      onForward={handleForward}
      onViewThread={(threadId) => setViewingThread(threadId)}
    />
  );

  return (
    <>
      <ThreePaneLayout
        sidebar={
        <Sidebar
          currentRecipient={recipient}
          recipients={recipients}
          domains={domains}
          selectedDomainId={selectedDomainId}
          onRecipientChange={(r) => {
            setRecipient(r);
            setSelectedMail(null);
          }}
          onDomainChange={async (domainId) => {
            setSelectedDomainId(domainId);
            const settings = await getUserSettings();
            await updateUserSettings({
              discord_webhook_url: settings.discord_webhook_url || "",
              selected_domain_id: domainId,
            });
            setRecipient("");
            setSelectedMail(null);
            const recipientRes = await getRecipients();
            setRecipients(recipientRes.recipients || []);
            await sync();
            await refresh();
          }}
          onCompose={handleCompose}
          onSettings={() => router.push("/settings")}
          onSignOut={signOut}
        />
        }
        list={
          <MailList
            mails={mails}
            selectedKey={selectedMail?.s3_key ?? null}
            onSelect={handleSelectMail}
            onToggleStar={toggleStar}
            loading={loading}
            page={page}
            totalPages={totalPages}
            onPageChange={setPage}
            onSync={sync}
          />
        }
        detail={detailPane}
        showDetail={!!selectedMail || !!viewingThread}
        onBackToList={() => {
          setSelectedMail(null);
          setViewingThread(null);
        }}
      />

      {showCompose && (
        <ComposeForm
          {...composeProps}
          onClose={() => setShowCompose(false)}
          onSent={() => {
            setShowCompose(false);
            refresh();
          }}
        />
      )}
    </>
  );
}
