export type ApiEnvelope<T> = {
  code: number;
  message: string;
  data: T;
};

export type UserProfile = {
  user_id: string;
  username: string;
  nickname: string;
  avatar: string;
  status: "online" | "offline" | string;
};

export type BackendMessage = {
  message_id: string;
  sender_id: string;
  receiver_id: string;
  content: string;
  msg_type?: string;
  message_type?: string;
  status: string;
  created_at?: string;
  read_at?: string;
};
