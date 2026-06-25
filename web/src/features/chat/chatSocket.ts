import { WS_BASE_URL } from "../../shared/lib/config";
import type { SendMessageInput } from "./chatTypes";

type WebSocketLike = {
  onclose: (() => void) | null;
  onerror: (() => void) | null;
  onmessage: ((event: { data: string }) => void) | null;
  onopen: (() => void) | null;
  close: () => void;
  send: (data: string) => void;
};

type WebSocketConstructor = new (url: string) => WebSocketLike;

type SocketStatus = "connecting" | "open" | "closed" | "error";

type AckPayload = {
  message_id: string;
  status: string;
};

type UnreadPayload = {
  count: number;
};

type SocketErrorPayload = {
  code?: number;
  message: string;
};

type ChatSocketCallbacks = {
  onAck?: (payload: AckPayload) => void;
  onMessage?: (payload: Record<string, unknown>) => void;
  onUnread?: (payload: UnreadPayload) => void;
  onError?: (payload: SocketErrorPayload) => void;
  onStatus?: (status: SocketStatus) => void;
};

type CreateChatSocketOptions = {
  token: string;
  baseUrl?: string;
  WebSocketImpl?: WebSocketConstructor;
  callbacks?: ChatSocketCallbacks;
};

type BackendSocketMessage = {
  event: string;
  payload?: unknown;
};

export function createChatSocket(options: CreateChatSocketOptions) {
  const callbacks = { ...options.callbacks };
  const WebSocketImpl =
    options.WebSocketImpl ?? (WebSocket as unknown as WebSocketConstructor);
  let ws: WebSocketLike | null = null;

  const connect = () => {
    if (ws) {
      return;
    }

    callbacks.onStatus?.("connecting");
    ws = new WebSocketImpl(buildSocketUrl(options.baseUrl ?? WS_BASE_URL, options.token));
    ws.onopen = () => callbacks.onStatus?.("open");
    ws.onclose = () => callbacks.onStatus?.("closed");
    ws.onerror = () => callbacks.onStatus?.("error");
    ws.onmessage = (event) => dispatchMessage(event.data, callbacks);
  };

  const sendEvent = (event: string, payload: unknown = {}) => {
    ws?.send(JSON.stringify({ event, payload }));
  };

  connect();

  return {
    connect,
    disconnect: () => {
      ws?.close();
      ws = null;
    },
    sendMessage: (payload: SendMessageInput) => {
      sendEvent("message", payload);
    },
    sendPing: () => {
      sendEvent("ping", {});
    },
    onAck: (callback: NonNullable<ChatSocketCallbacks["onAck"]>) => {
      callbacks.onAck = callback;
    },
    onMessage: (callback: NonNullable<ChatSocketCallbacks["onMessage"]>) => {
      callbacks.onMessage = callback;
    },
    onUnread: (callback: NonNullable<ChatSocketCallbacks["onUnread"]>) => {
      callbacks.onUnread = callback;
    },
    onError: (callback: NonNullable<ChatSocketCallbacks["onError"]>) => {
      callbacks.onError = callback;
    },
    onStatus: (callback: NonNullable<ChatSocketCallbacks["onStatus"]>) => {
      callbacks.onStatus = callback;
    }
  };
}

function buildSocketUrl(baseUrl: string, token: string) {
  const url = new URL(baseUrl);
  url.searchParams.set("token", token);
  return url.toString();
}

function dispatchMessage(raw: string, callbacks: ChatSocketCallbacks) {
  let message: BackendSocketMessage;

  try {
    message = JSON.parse(raw) as BackendSocketMessage;
  } catch {
    callbacks.onError?.({ message: "Invalid WebSocket message" });
    return;
  }

  switch (message.event) {
    case "ack":
      callbacks.onAck?.(message.payload as AckPayload);
      break;
    case "message":
      callbacks.onMessage?.(message.payload as Record<string, unknown>);
      break;
    case "unread":
      callbacks.onUnread?.(message.payload as UnreadPayload);
      break;
    case "error":
      callbacks.onError?.(message.payload as SocketErrorPayload);
      break;
    case "pong":
      break;
    default:
      callbacks.onError?.({ message: `Unknown WebSocket event: ${message.event}` });
  }
}
