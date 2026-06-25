import { beforeEach, describe, expect, it, vi } from "vitest";

import { createChatSocket } from "./chatSocket";

class FakeWebSocket {
  static lastUrl = "";
  static lastSent = "";
  static instance: FakeWebSocket | null = null;

  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onopen: (() => void) | null = null;

  constructor(url: string) {
    FakeWebSocket.lastUrl = url;
    FakeWebSocket.instance = this;
  }

  send(data: string) {
    FakeWebSocket.lastSent = data;
  }

  close() {
    this.onclose?.();
  }

  emit(event: string, payload: unknown) {
    this.onmessage?.({ data: JSON.stringify({ event, payload }) });
  }
}

describe("createChatSocket", () => {
  beforeEach(() => {
    FakeWebSocket.lastUrl = "";
    FakeWebSocket.lastSent = "";
    FakeWebSocket.instance = null;
  });

  it("connects with a query token and sends backend message events", () => {
    const socket = createChatSocket({
      token: "jwt",
      baseUrl: "ws://localhost:8080/api/v1/ws",
      WebSocketImpl: FakeWebSocket
    });

    expect(FakeWebSocket.lastUrl).toBe(
      "ws://localhost:8080/api/v1/ws?token=jwt"
    );

    socket.sendMessage({
      receiver_id: "u2",
      content: "hello",
      message_type: "text"
    });

    expect(FakeWebSocket.lastSent).toContain('"event":"message"');
    expect(FakeWebSocket.lastSent).toContain('"message_type":"text"');
  });

  it("dispatches backend socket events", () => {
    const onAck = vi.fn();
    const onMessage = vi.fn();
    const onUnread = vi.fn();
    const onError = vi.fn();

    createChatSocket({
      token: "jwt",
      baseUrl: "ws://localhost:8080/api/v1/ws",
      WebSocketImpl: FakeWebSocket,
      callbacks: {
        onAck,
        onMessage,
        onUnread,
        onError
      }
    });

    FakeWebSocket.instance?.emit("ack", { message_id: "m1", status: "sent" });
    FakeWebSocket.instance?.emit("message", {
      message_id: "m2",
      sender_id: "u2",
      receiver_id: "u1",
      content: "hi",
      message_type: "text"
    });
    FakeWebSocket.instance?.emit("unread", { count: 3 });
    FakeWebSocket.instance?.emit("error", { code: 3001, message: "failed" });

    expect(onAck).toHaveBeenCalledWith({ message_id: "m1", status: "sent" });
    expect(onMessage).toHaveBeenCalledWith(
      expect.objectContaining({ message_id: "m2", content: "hi" })
    );
    expect(onUnread).toHaveBeenCalledWith({ count: 3 });
    expect(onError).toHaveBeenCalledWith({ code: 3001, message: "failed" });
  });
});
