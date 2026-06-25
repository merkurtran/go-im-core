import {
  CheckCheck,
  MessageCircle,
  Plus,
} from "lucide-react";
import { FormEvent, useState } from "react";
import { useAuthStore } from "../../auth/authStore";
import { useChatStore } from "../chatStore";

export function ConversationSidebar() {
  const user = useAuthStore((s) => s.user);
  const contacts = useChatStore((s) => s.contacts);
  const selectedContactId = useChatStore((s) => s.selectedContactId);
  const selectContact = useChatStore((s) => s.selectContact);
  const upsertContact = useChatStore((s) => s.upsertContact);

  const [contactIdInput, setContactIdInput] = useState("");
  const [contactNameInput, setContactNameInput] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);

  const handleAddContact = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const userId = contactIdInput.trim();
    if (!userId) return;
    upsertContact({
      user_id: userId,
      display_name: contactNameInput.trim() || userId,
      avatar: "",
    });
    selectContact(userId);
    setContactIdInput("");
    setContactNameInput("");
    setShowAddForm(false);
  };

  return (
    <aside className="conversation-sidebar" aria-label="会话列表">
      <div className="sidebar-header">
        <h2>Messages</h2>
        <button
          className="icon-btn"
          type="button"
          aria-label="添加联系人"
          onClick={() => setShowAddForm((v) => !v)}
        >
          <Plus size={18} />
        </button>
      </div>

      {showAddForm ? (
        <form className="add-contact-form compact" onSubmit={handleAddContact}>
          <label>
            对方 user_id
            <input
              value={contactIdInput}
              onChange={(e) => setContactIdInput(e.target.value)}
              placeholder="输入 user_id"
            />
          </label>
          <label>
            备注名
            <input
              value={contactNameInput}
              onChange={(e) => setContactNameInput(e.target.value)}
              placeholder="可选"
            />
          </label>
          <div className="form-actions">
            <button className="button button-secondary" type="submit">
              添加
            </button>
            <button
              className="button ghost-button"
              type="button"
              onClick={() => setShowAddForm(false)}
            >
              取消
            </button>
          </div>
        </form>
      ) : null}

      <div className="contact-list">
        {contacts.length === 0 ? (
          <div className="empty-list">
            <MessageCircle size={22} />
            <p>添加一个联系人，用 user_id 开始第一段会话。</p>
          </div>
        ) : (
          contacts.map((contact) => {
            const isActive = contact.user_id === selectedContactId;
            const isOnline = contact.status === "online";
            const isTyping = contact.is_typing;
            return (
              <button
                className={`contact-item ${isActive ? "is-active" : ""}`}
                key={contact.user_id}
                type="button"
                onClick={() => selectContact(contact.user_id)}
              >
                <div className="contact-avatar-wrap">
                  <div className="contact-avatar">
                    {contact.display_name.trim().slice(0, 1).toUpperCase()}
                  </div>
                  {isOnline ? <span className="online-dot" /> : null}
                </div>
                <div className="contact-info">
                  <div className="contact-row">
                    <strong>{contact.display_name}</strong>
                    <span className="contact-time">
                      {contact.last_message_time || ""}
                    </span>
                  </div>
                  <div className="contact-row">
                    <span className="contact-preview">
                      {isTyping ? (
                        <span className="typing-text">...is typing</span>
                      ) : (
                        contact.last_message || ""
                      )}
                    </span>
                    <span className="contact-meta">
                      {contact.unread_count ? (
                        <span className="unread-badge">
                          {contact.unread_count}
                        </span>
                      ) : (
                        <CheckCheck size={14} className="read-check" />
                      )}
                    </span>
                  </div>
                </div>
              </button>
            );
          })
        )}
      </div>

      <div className="current-user-mini">
        <div className="contact-avatar-wrap">
          <div className="contact-avatar">
            {(user?.nickname || user?.username || "U")
              .trim()
              .slice(0, 1)
              .toUpperCase()}
          </div>
          <span className="online-dot" />
        </div>
        <div className="current-user-info">
          <strong>{user?.nickname || user?.username || "未登录"}</strong>
          <span>{user?.user_id || "等待认证"}</span>
        </div>
      </div>
    </aside>
  );
}
