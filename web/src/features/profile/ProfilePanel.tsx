import { useMutation } from "@tanstack/react-query";
import { Save } from "lucide-react";
import { FormEvent, useEffect, useState } from "react";

import { updateProfile } from "../../api/users";
import { useAuthStore } from "../auth/authStore";

export function ProfilePanel() {
  const token = useAuthStore((state) => state.token);
  const user = useAuthStore((state) => state.user);
  const setUser = useAuthStore((state) => state.setUser);
  const [nickname, setNickname] = useState(user?.nickname ?? "");
  const [avatar, setAvatar] = useState(user?.avatar ?? "");

  useEffect(() => {
    setNickname(user?.nickname ?? "");
    setAvatar(user?.avatar ?? "");
  }, [user?.avatar, user?.nickname]);

  const mutation = useMutation({
    mutationFn: async () => {
      if (!token || !user) {
        return null;
      }
      const updated = await updateProfile(token, { nickname, avatar });
      return updated ?? { ...user, nickname, avatar };
    },
    onSuccess: (updated) => {
      if (updated) {
        setUser(updated);
      }
    }
  });

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    mutation.mutate();
  };

  if (!user) {
    return (
      <section className="side-section">
        <h2>个人资料</h2>
        <p className="subtle">登录后可编辑昵称和头像。</p>
      </section>
    );
  }

  return (
    <section className="side-section">
      <div className="section-title-row">
        <h2>个人资料</h2>
        <span className="status-pill">{user.status || "offline"}</span>
      </div>
      <form className="profile-form" onSubmit={handleSubmit}>
        <label>
          昵称
          <input
            value={nickname}
            maxLength={20}
            onChange={(event) => setNickname(event.target.value)}
          />
        </label>
        <label>
          头像 URL
          <input
            value={avatar}
            onChange={(event) => setAvatar(event.target.value)}
          />
        </label>
        {mutation.isError ? (
          <p className="form-error">资料保存失败，请稍后重试。</p>
        ) : null}
        <button className="button button-secondary" type="submit">
          <Save size={16} />
          保存资料
        </button>
      </form>
    </section>
  );
}
