import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { useAuthStore } from "../auth/authStore";
import { ChatPage } from "./ChatPage";
import { useChatStore } from "./chatStore";

class FakeWebSocket {
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onopen: (() => void) | null = null;
  constructor(_url: string) {}
  close() {
    this.onclose?.();
  }
  send(_data: string) {}
}

function renderChatPage() {
  const queryClient = new QueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <ChatPage />
    </QueryClientProvider>
  );
}

describe("ChatPage", () => {
  beforeEach(() => {
    vi.stubGlobal("WebSocket", FakeWebSocket);
    useAuthStore.setState({
      token: "jwt",
      user: {
        user_id: "u1",
        username: "alice",
        nickname: "Alice",
        avatar: "",
        status: "offline"
      }
    });
    useChatStore.setState({
      contacts: [],
      selectedContactId: null,
      messagesByContact: {}
    });
  });

  it("renders the empty contact state and disables sending without a contact", () => {
    renderChatPage();

    expect(
      screen.getByText(/添加一个联系人/i)
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /发送/i })).toBeDisabled();
  });
});
