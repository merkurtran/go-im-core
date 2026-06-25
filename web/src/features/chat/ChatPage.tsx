import { LeftNavBar } from "./components/LeftNavBar";
import { ConversationSidebar } from "./components/ConversationSidebar";
import { ChatWindow } from "./components/ChatWindow";

export function ChatPage() {
  return (
    <main className="chat-layout">
      <LeftNavBar />
      <ConversationSidebar />
      <ChatWindow />
    </main>
  );
}
