import {
  Bookmark,
  Bell,
  MoreHorizontal,
  Paperclip,
  Play,
  Send,
  UserRound,
  RefreshCw,
  MessageCircle,
  Eye,
  Building,
  Thermometer,
} from "lucide-react";
import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { getConversation } from "../../../api/messages";
import type { BackendMessage } from "../../../api/types";
import { useAuthStore } from "../../auth/authStore";
import { useChatStore } from "../chatStore";
import { getMockMessages } from "../chatStore";
import type { ChatMessage } from "../chatTypes";
import { createChatSocket } from "../chatSocket";

export function ChatWindow() {
  const token = useAuthStore((s) => s.token);
  const user = useAuthStore((s) => s.user);
  const selectedContactId = useChatStore((s) => s.selectedContactId);
  const contacts = useChatStore((s) => s.contacts);
  const messagesByContact = useChatStore((s) => s.messagesByContact);
  const addOptimisticMessage = useChatStore((s) => s.addOptimisticMessage);
  const confirmMessage = useChatStore((s) => s.confirmMessage);
  const receiveMessage = useChatStore((s) => s.receiveMessage);
  const markMessageFailed = useChatStore((s) => s.markMessageFailed);
  const upsertContact = useChatStore((s) => s.upsertContact);
  const setContactMessages = useChatStore((s) => s.setContactMessages);

  const [draft, setDraft] = useState("");
  const [socketStatus, setSocketStatus] = useState("closed");
  const [socketError, setSocketError] = useState("");
  const socketRef = useRef<ReturnType<typeof createChatSocket> | null>(null);
  const pendingClientIds = useRef<string[]>([]);
  const bottomRef = useRef<HTMLDivElement | null>(null);

  const selectedContact = contacts.find(
    (c) => c.user_id === selectedContactId
  );

  const historyQuery = useQuery({
    queryKey: ["conversation", selectedContactId],
    queryFn: () => getConversation(token ?? "", selectedContactId ?? ""),
    enabled: Boolean(token && selectedContactId),
  });

  // 加载 mock 消息
  useEffect(() => {
    if (selectedContactId && messagesByContact[selectedContactId]?.length === 0) {
      const mocks = getMockMessages(selectedContactId);
      if (mocks.length > 0) {
        setContactMessages(selectedContactId, mocks);
      }
    }
  }, [selectedContactId, messagesByContact, setContactMessages]);

  useEffect(() => {
    if (!token) return;
    const socket = createChatSocket({
      token,
      callbacks: {
        onStatus: setSocketStatus,
        onAck: (payload) => {
          const clientId = pendingClientIds.current.shift();
          if (clientId) {
            confirmMessage(clientId, payload.message_id);
          }
        },
        onMessage: (payload) => {
          const message = toChatMessage(payload as BackendMessage, user?.user_id);
          upsertContact({
            user_id: message.sender_id,
            display_name: message.sender_id,
            avatar: "",
          });
          receiveMessage(message);
        },
        onUnread: (payload) => {
          setSocketError(
            payload.count > 0 ? `有 ${payload.count} 条未读消息` : ""
          );
        },
        onError: (payload) => {
          setSocketError(payload.message);
        },
      },
    });
    socketRef.current = socket;
    return () => {
      socket.disconnect();
      socketRef.current = null;
    };
  }, [confirmMessage, receiveMessage, token, upsertContact, user?.user_id]);

  const historyMessages = useMemo(() => {
    return (historyQuery.data ?? []).map((message) =>
      toChatMessage(message, user?.user_id)
    );
  }, [historyQuery.data, user?.user_id]);

  const localMessages = selectedContactId
    ? messagesByContact[selectedContactId] ?? []
    : [];

  const messages = useMemo(() => {
    const merged = [...historyMessages, ...localMessages];
    // 按时间排序
    merged.sort(
      (a, b) =>
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    );
    return merged;
  }, [historyMessages, localMessages]);

  // 自动滚动到底部
  useEffect(() => {
    if (
      bottomRef.current &&
      typeof bottomRef.current.scrollIntoView === "function"
    ) {
      bottomRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages.length]);

  const handleSend = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!selectedContact || !user || !draft.trim()) return;

    const clientId = createClientId();
    const message: ChatMessage = {
      client_id: clientId,
      sender_id: user.user_id,
      receiver_id: selectedContact.user_id,
      content: draft.trim(),
      message_type: "text",
      status: "pending",
      created_at: new Date().toISOString(),
    };

    addOptimisticMessage(message);
    pendingClientIds.current.push(clientId);
    setDraft("");

    try {
      socketRef.current?.sendMessage({
        receiver_id: selectedContact.user_id,
        content: message.content,
        message_type: "text",
      });
    } catch {
      markMessageFailed(clientId);
    }
  };

  return (
    <section className="chat-window" aria-label="聊天窗口">
      <ChatHeader
        contact={selectedContact}
        socketStatus={socketStatus}
        socketError={socketError}
      />

      <div className="message-timeline">
        {!selectedContact ? (
          <div className="empty-chat">
            <UserRound size={32} />
            <h2>选择一个联系人</h2>
            <p>从左侧列表选择联系人开始聊天。</p>
          </div>
        ) : historyQuery.isLoading ? (
          <div className="empty-chat">
            <RefreshCw size={28} />
            <p>正在加载历史消息...</p>
          </div>
        ) : messages.length === 0 ? (
          <div className="empty-chat">
            <MessageCircle size={30} />
            <p>还没有消息，发送一句话开始聊天。</p>
          </div>
        ) : (
          messages.map((message) => (
            <MessageBubble
              key={message.message_id ?? message.client_id}
              message={message}
              isMine={message.sender_id === user?.user_id}
            />
          ))
        )}
        <div ref={bottomRef} />
      </div>

      <form className="composer" onSubmit={handleSend}>
        <button
          className="attach-btn"
          type="button"
          aria-label="附件"
          disabled={!selectedContact}
        >
          <Paperclip size={20} />
        </button>
        <input
          className="composer-input"
          aria-label="消息内容"
          disabled={!selectedContact}
          placeholder={
            selectedContact ? "Write a Message" : "先选择一个联系人"
          }
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
        />
        <button
          className="send-circle-btn"
          disabled={!selectedContact || !draft.trim()}
          type="submit"
          aria-label="发送"
        >
          <Send size={18} />
        </button>
      </form>
    </section>
  );
}

function ChatHeader({
  contact,
  socketStatus,
  socketError,
}: {
  contact?: { display_name: string; status?: string; avatar?: string };
  socketStatus: string;
  socketError: string;
}) {
  if (!contact) {
    return (
      <header className="chat-header">
        <div className="chat-header-info">
          <h1>选择联系人</h1>
        </div>
      </header>
    );
  }

  const isOnline = contact.status === "online" || socketStatus === "open";

  return (
    <header className="chat-header">
      <div className="chat-header-info">
        <div className="header-avatar">
          {contact.avatar ? (
            <img src={contact.avatar} alt="" />
          ) : (
            <span>{contact.display_name.trim().slice(0, 1).toUpperCase()}</span>
          )}
          {isOnline ? <span className="online-dot" /> : null}
        </div>
        <div>
          <h1>{contact.display_name}</h1>
          <span className="header-status">
            {isOnline ? "Online" : "Offline"}
          </span>
        </div>
      </div>
      <div className="chat-header-actions">
        <button className="icon-btn" type="button" aria-label="收藏">
          <Bookmark size={20} />
        </button>
        <button className="icon-btn" type="button" aria-label="通知">
          <Bell size={20} />
        </button>
        <button className="icon-btn" type="button" aria-label="更多">
          <MoreHorizontal size={20} />
        </button>
      </div>
      {socketError ? (
        <p className="inline-notice header-notice">{socketError}</p>
      ) : null}
    </header>
  );
}

function MessageBubble({
  message,
  isMine,
}: {
  message: ChatMessage;
  isMine: boolean;
}) {
  const time = formatTime(message.created_at);

  if (message.message_type === "system") {
    return (
      <div className="message-row system">
        <div className="message-bubble system-bubble">
          <p>{message.content}</p>
          <span className="msg-time">{time}</span>
        </div>
      </div>
    );
  }

  if (message.message_type === "voice") {
    return (
      <div className={`message-row ${isMine ? "mine" : ""}`}>
        <div className="message-bubble voice-bubble">
          <button className="voice-play" type="button" aria-label="播放">
            <Play size={16} fill="currentColor" />
          </button>
          <div className="voice-wave">
            {Array.from({ length: 24 }).map((_, i) => (
              <span
                key={i}
                style={{
                  height: `${Math.max(4, Math.random() * 18 + 4)}px`,
                }}
              />
            ))}
          </div>
          <span className="voice-duration">
            {message.extra?.duration || "0:00"}
          </span>
          <span className="msg-time">{time}</span>
        </div>
      </div>
    );
  }

  if (message.message_type === "card") {
    return (
      <div className={`message-row ${isMine ? "mine" : ""}`}>
        <div className="message-bubble card-bubble">
          <p className="card-text">{message.content}</p>
          {message.extra?.cards ? (
            <div className="card-grid">
              {message.extra.cards.map((card, idx) => (
                <div className="card-item" key={idx}>
                  <div className="card-icon">
                    {card.icon === "eye" && <Eye size={16} />}
                    {card.icon === "building" && <Building size={16} />}
                    {card.icon === "thermometer" && <Thermometer size={16} />}
                  </div>
                  <span className="card-label">{card.label}</span>
                  <strong className="card-value">{card.value}</strong>
                </div>
              ))}
            </div>
          ) : null}
          {message.extra?.images ? (
            <div className="image-grid small">
              {message.extra.images.slice(0, 3).map((_, idx) => (
                <div className="image-placeholder" key={idx}>
                  <span>Photo {idx + 1}</span>
                </div>
              ))}
            </div>
          ) : null}
          <span className="msg-time">{time}</span>
        </div>
      </div>
    );
  }

  if (message.message_type === "image") {
    return (
      <div className={`message-row ${isMine ? "mine" : ""}`}>
        <div className="message-bubble image-bubble">
          {message.content ? <p className="image-caption">{message.content}</p> : null}
          {message.extra?.images ? (
            <div className="image-grid">
              {message.extra.images.slice(0, 3).map((_, idx) => (
                <div className="image-placeholder" key={idx}>
                  <span>{idx === 2 && message.extra!.images!.length > 3 ? "3+" : `Photo ${idx + 1}`}</span>
                </div>
              ))}
            </div>
          ) : null}
          <span className="msg-time">{time}</span>
        </div>
      </div>
    );
  }

  // text / file default
  return (
    <div className={`message-row ${isMine ? "mine" : ""}`}>
      <div className={`message-bubble ${isMine ? "mine" : ""}`}>
        <p>{message.content}</p>
        <span className="msg-time">{time}</span>
      </div>
    </div>
  );
}

function formatTime(iso: string) {
  try {
    const d = new Date(iso);
    return d.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    });
  } catch {
    return "";
  }
}

function toChatMessage(
  message: BackendMessage,
  currentUserId?: string
): ChatMessage {
  return {
    client_id: message.message_id || createClientId(),
    message_id: message.message_id,
    sender_id: message.sender_id,
    receiver_id: message.receiver_id,
    content: message.content,
    message_type: (message.message_type ?? message.msg_type ?? "text") as
      | "text"
      | "image"
      | "file",
    status:
      message.sender_id === currentUserId
        ? "sent"
        : ((message.status as ChatMessage["status"]) || "received"),
    created_at: message.created_at ?? new Date().toISOString(),
  };
}

function createClientId() {
  if (typeof crypto !== "undefined" && "randomUUID" in crypto) {
    return crypto.randomUUID();
  }
  return `client-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}
