export type ChatContact = {
  user_id: string;
  display_name: string;
  avatar?: string;
  status?: "online" | "offline" | "away";
  is_typing?: boolean;
  last_message?: string;
  last_message_time?: string;
  unread_count?: number;
};

export type ChatMessageStatus = "pending" | "sent" | "received" | "read" | "failed";

export type ChatMessage = {
  client_id: string;
  message_id?: string;
  sender_id: string;
  sender_name?: string;
  sender_avatar?: string;
  receiver_id: string;
  content: string;
  message_type: "text" | "image" | "file" | "card" | "voice" | "system";
  status: ChatMessageStatus;
  created_at: string;
  extra?: {
    images?: string[];
    duration?: string;
    cards?: { label: string; value: string; icon?: string }[];
  };
};

export type SendMessageInput = {
  receiver_id: string;
  content: string;
  message_type: "text";
};
