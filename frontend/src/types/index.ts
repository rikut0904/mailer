export interface ParsedMail {
  s3_key: string;
  message_id: string;
  from: string;
  to: string;
  subject: string;
  body: string;
  html_body?: string;
  date: string;
  attachments?: Attachment[];
  is_read: boolean;
  is_starred: boolean;
  thread_id?: string;
}

export interface Attachment {
  filename: string;
  content_type: string;
  size: number;
}

export interface MailListResponse {
  mails: ParsedMail[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface ThreadGroup {
  parent_uuid: string;
  group_name: string;
}

export interface ThreadMessage {
  type: "received" | "sent";
  subject: string;
  from: string;
  to: string;
  body: string;
  date: string;
  s3_key?: string;
  management_code?: string;
  is_read?: boolean;
  is_starred?: boolean;
}

export interface ThreadResponse {
  thread_id: string;
  group_name: string;
  messages: ThreadMessage[];
}

export interface SendRequest {
  to: string[];
  subject: string;
  body: string;
  html_body?: string;
  thread_id?: string;
  reply_code?: string;
  send_type: "new" | "reply" | "forward";
  from_address?: string;
}

export interface SendResponse {
  thread_id: string;
  management_codes: string[];
}

export interface UserSettings {
  discord_webhook_url: string;
  selected_domain_id: string;
}

export interface S3Domain {
  id: string;
  name: string;
  bucket: string;
  region: string;
  access_key_id: string;
  secret_key: string;
  endpoint: string;
}
