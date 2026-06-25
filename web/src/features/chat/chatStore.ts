import { create } from "zustand";

import type { ChatContact, ChatMessage } from "./chatTypes";

const CONTACTS_STORAGE_KEY = "go-im-core.contacts";

function getMockContacts(): ChatContact[] {
  return [
    {
      user_id: "john-doe",
      display_name: "John Doe",
      status: "online",
      last_message: "How's going with your property search?",
      last_message_time: "4:45 pm",
      unread_count: 0
    },
    {
      user_id: "travis-barker",
      display_name: "Travis Barker",
      status: "online",
      is_typing: true,
      last_message: "...is typing",
      last_message_time: "5:38 pm",
      unread_count: 0
    },
    {
      user_id: "kate-rose",
      display_name: "Kate Rose",
      status: "offline",
      last_message: "Looking forward to discussing real estate...",
      last_message_time: "5:04 pm",
      unread_count: 0
    },
    {
      user_id: "robert-parker",
      display_name: "Robert Parker",
      status: "online",
      last_message: "That's fantastic news about the new listi...",
      last_message_time: "4:22 pm",
      unread_count: 1
    },
    {
      user_id: "emily-johnson",
      display_name: "Emily Johnson",
      status: "offline",
      last_message: "Take a look at my recent real estate post...",
      last_message_time: "3:59 pm",
      unread_count: 0
    },
    {
      user_id: "sophia-brown",
      display_name: "Sophia Brown",
      status: "online",
      last_message: "Discover amazing properties on my page!",
      last_message_time: "3:24 pm",
      unread_count: 0
    },
    {
      user_id: "frederick-k",
      display_name: "Frederick K.",
      status: "offline",
      last_message: "Are you considering pest control options...",
      last_message_time: "1:06 pm",
      unread_count: 0
    },
    {
      user_id: "alexander-jameson",
      display_name: "Alexander Jameson",
      status: "online",
      last_message: "What day and time works best for you to...",
      last_message_time: "12:57 am",
      unread_count: 0
    },
    {
      user_id: "tom-hardy",
      display_name: "Tom Hardy",
      status: "away",
      last_message: "This property has such a unique design vi...",
      last_message_time: "12:43 pm",
      unread_count: 0
    },
    {
      user_id: "vivienne-marigold",
      display_name: "Vivienne Marigold",
      status: "offline",
      last_message: "",
      last_message_time: "12:17 pm",
      unread_count: 0
    }
  ];
}

export function getMockMessages(contactId: string): ChatMessage[] {
  if (contactId !== "alexander-jameson") return [];
  const now = new Date();
  const base = new Date(now);
  base.setHours(base.getHours() - 2);
  return [
    {
      client_id: "mock-1",
      sender_id: "alexander-jameson",
      sender_name: "Alexander Jameson",
      receiver_id: "me",
      content: "I'm a manager that's here to help",
      message_type: "system",
      status: "received",
      created_at: new Date(base.getTime() - 1000 * 60 * 120).toISOString()
    },
    {
      client_id: "mock-2",
      sender_id: "alexander-jameson",
      sender_name: "Alexander Jameson",
      receiver_id: "me",
      content: "",
      message_type: "voice",
      status: "received",
      created_at: new Date(base.getTime() - 1000 * 60 * 116).toISOString(),
      extra: { duration: "2:19" }
    },
    {
      client_id: "mock-3",
      sender_id: "alexander-jameson",
      sender_name: "Alexander Jameson",
      receiver_id: "me",
      content: "This is a modern townhouse in a quiet neighborhood. It has 3 bedrooms, 2 bathrooms, a spacious kitchen and living room 🏠",
      message_type: "card",
      status: "received",
      created_at: new Date(base.getTime() - 1000 * 60 * 115).toISOString(),
      extra: {
        cards: [
          { label: "Daily Visitors", value: "2,429", icon: "eye" },
          { label: "Building Age", value: "3Y", icon: "building" },
          { label: "Temperature", value: "28°F", icon: "thermometer" }
        ],
        images: ["img1", "img2", "img3"]
      }
    },
    {
      client_id: "mock-4",
      sender_id: "alexander-jameson",
      sender_name: "Alexander Jameson",
      receiver_id: "me",
      content: "Here are the photos:",
      message_type: "image",
      status: "received",
      created_at: new Date(base.getTime() - 1000 * 60 * 78).toISOString(),
      extra: { images: ["img1", "img2", "img3"] }
    },
    {
      client_id: "mock-5",
      sender_id: "me",
      sender_name: "Me",
      receiver_id: "alexander-jameson",
      content: "Looks good 👍, I want to sign up for a viewing",
      message_type: "text",
      status: "sent",
      created_at: new Date(base.getTime() - 1000 * 60 * 32).toISOString()
    },
    {
      client_id: "mock-6",
      sender_id: "alexander-jameson",
      sender_name: "Alexander Jameson",
      receiver_id: "me",
      content: "What day and time works best for you to come by for the viewing? Let me know, and I'll confirm the appointment right away.",
      message_type: "text",
      status: "received",
      created_at: new Date(base.getTime() - 1000 * 60).toISOString()
    }
  ];
}

type ChatState = {
  contacts: ChatContact[];
  selectedContactId: string | null;
  messagesByContact: Record<string, ChatMessage[]>;
  upsertContact: (contact: ChatContact) => void;
  removeContact: (userId: string) => void;
  selectContact: (userId: string | null) => void;
  addOptimisticMessage: (message: ChatMessage) => void;
  confirmMessage: (clientId: string, messageId: string) => void;
  receiveMessage: (message: ChatMessage) => void;
  markMessageFailed: (clientId: string) => void;
  setContactMessages: (contactId: string, messages: ChatMessage[]) => void;
};

export const useChatStore = create<ChatState>((set, get) => ({
  contacts: readStoredContacts(),
  selectedContactId: null,
  messagesByContact: {},
  upsertContact: (contact) => {
    const contacts = mergeContact(get().contacts, contact);
    writeStoredContacts(contacts);
    set({ contacts });
  },
  removeContact: (userId) => {
    const contacts = get().contacts.filter((contact) => contact.user_id !== userId);
    const { [userId]: _removed, ...messagesByContact } = get().messagesByContact;
    writeStoredContacts(contacts);
    set({
      contacts,
      messagesByContact,
      selectedContactId:
        get().selectedContactId === userId ? null : get().selectedContactId
    });
  },
  selectContact: (userId) => {
    set({ selectedContactId: userId });
  },
  addOptimisticMessage: (message) => {
    appendMessage(set, get, message.receiver_id, message);
  },
  confirmMessage: (clientId, messageId) => {
    updateMessage(set, get, clientId, (message) => ({
      ...message,
      message_id: messageId,
      status: "sent"
    }));
  },
  receiveMessage: (message) => {
    const contactId = getConversationId(message, get().contacts);
    appendMessage(set, get, contactId, message);
  },
  markMessageFailed: (clientId) => {
    updateMessage(set, get, clientId, (message) => ({
      ...message,
      status: "failed"
    }));
  },
  setContactMessages: (contactId, messages) => {
    set({
      messagesByContact: {
        ...get().messagesByContact,
        [contactId]: messages
      }
    });
  }
}));

function mergeContact(contacts: ChatContact[], next: ChatContact) {
  const index = contacts.findIndex((contact) => contact.user_id === next.user_id);
  if (index === -1) {
    return [next, ...contacts];
  }

  const merged = [...contacts];
  merged[index] = { ...merged[index], ...next };
  return merged;
}

function appendMessage(
  set: (state: Partial<ChatState>) => void,
  get: () => ChatState,
  contactId: string,
  message: ChatMessage
) {
  const current = get().messagesByContact[contactId] ?? [];
  set({
    messagesByContact: {
      ...get().messagesByContact,
      [contactId]: [...current, message]
    }
  });
}

function updateMessage(
  set: (state: Partial<ChatState>) => void,
  get: () => ChatState,
  clientId: string,
  update: (message: ChatMessage) => ChatMessage
) {
  const messagesByContact = Object.fromEntries(
    Object.entries(get().messagesByContact).map(([contactId, messages]) => [
      contactId,
      messages.map((message) =>
        message.client_id === clientId ? update(message) : message
      )
    ])
  );

  set({ messagesByContact });
}

function getConversationId(message: ChatMessage, contacts: ChatContact[]) {
  const senderIsKnownContact = contacts.some(
    (contact) => contact.user_id === message.sender_id
  );

  return senderIsKnownContact ? message.sender_id : message.receiver_id;
}

function readStoredContacts(): ChatContact[] {
  if (typeof localStorage === "undefined") {
    return getMockContacts();
  }

  const raw = localStorage.getItem(CONTACTS_STORAGE_KEY);
  if (!raw) {
    return getMockContacts();
  }

  try {
    const contacts = JSON.parse(raw) as ChatContact[];
    return Array.isArray(contacts) && contacts.length > 0 ? contacts : getMockContacts();
  } catch {
    return getMockContacts();
  }
}

function writeStoredContacts(contacts: ChatContact[]) {
  if (typeof localStorage !== "undefined") {
    localStorage.setItem(CONTACTS_STORAGE_KEY, JSON.stringify(contacts));
  }
}
