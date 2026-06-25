import { beforeEach, describe, expect, it } from "vitest";

import { useChatStore } from "./chatStore";
import type { ChatMessage } from "./chatTypes";

const contact = { user_id: "u2", display_name: "Bob", avatar: "" };

function message(overrides: Partial<ChatMessage> = {}): ChatMessage {
  return {
    client_id: "c1",
    sender_id: "u1",
    receiver_id: "u2",
    content: "hello",
    message_type: "text",
    status: "pending",
    created_at: "2026-06-25T07:30:00.000Z",
    ...overrides
  };
}

describe("useChatStore", () => {
  beforeEach(() => {
    localStorage.clear();
    useChatStore.setState({
      contacts: [],
      selectedContactId: null,
      messagesByContact: {}
    });
  });

  it("upserts and selects local contacts", () => {
    useChatStore.getState().upsertContact(contact);

    expect(useChatStore.getState().contacts).toHaveLength(1);

    useChatStore.getState().selectContact("u2");

    expect(useChatStore.getState().selectedContactId).toBe("u2");
    expect(localStorage.getItem("go-im-core.contacts")).toContain("Bob");
  });

  it("updates optimistic messages after ack", () => {
    useChatStore.getState().addOptimisticMessage(message());
    useChatStore.getState().confirmMessage("c1", "m1");

    const [stored] = useChatStore.getState().messagesByContact.u2;

    expect(stored.message_id).toBe("m1");
    expect(stored.status).toBe("sent");
  });

  it("stores incoming messages in the sender conversation", () => {
    useChatStore.getState().upsertContact(contact);
    useChatStore.getState().receiveMessage(
      message({
        client_id: "server-m2",
        message_id: "m2",
        sender_id: "u2",
        receiver_id: "u1",
        content: "hi back",
        status: "received"
      })
    );

    const [stored] = useChatStore.getState().messagesByContact.u2;

    expect(stored.message_id).toBe("m2");
    expect(stored.content).toBe("hi back");
  });

  it("marks optimistic messages as failed", () => {
    useChatStore.getState().addOptimisticMessage(message({ client_id: "c3" }));
    useChatStore.getState().markMessageFailed("c3");

    const [stored] = useChatStore.getState().messagesByContact.u2;

    expect(stored.status).toBe("failed");
  });
});
