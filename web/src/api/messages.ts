import type { SendMessageInput } from "../features/chat/chatTypes";
import { apiRequest } from "./client";
import type { BackendMessage } from "./types";

export function getConversation(
  token: string,
  userId: string,
  limit = 50,
  offset = 0
) {
  return apiRequest<BackendMessage[]>("/messages", {
    token,
    query: {
      user_id: userId,
      limit,
      offset
    }
  });
}

export function sendRestMessage(token: string, input: SendMessageInput) {
  return apiRequest<{ message_id: string }>("/messages", {
    method: "POST",
    token,
    body: input
  });
}

export function markConversationRead(token: string, otherUserId: string) {
  return apiRequest<null>("/messages/read", {
    method: "PATCH",
    token,
    body: {
      other_user_id: otherUserId
    }
  });
}
