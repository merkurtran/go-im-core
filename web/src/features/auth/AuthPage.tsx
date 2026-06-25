import { useMutation } from "@tanstack/react-query";
import { MessageCircle } from "lucide-react";
import { FormEvent, useMemo, useState } from "react";

import { ApiError } from "../../api/client";
import { login, register } from "../../api/auth";
import { useAuthStore } from "./authStore";

type AuthMode = "login" | "register";

export function AuthPage() {
  const setSession = useAuthStore((state) => state.setSession);
  const [mode, setMode] = useState<AuthMode>("login");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [nickname, setNickname] = useState("");

  const authMutation = useMutation({
    mutationFn: async () => {
      if (mode === "login") {
        return login({ username, password });
      }

      const response = await register({
        username,
        password,
        nickname: nickname || username
      });

      return {
        token: response.token,
        user: {
          user_id: response.user_id,
          username,
          nickname: nickname || username,
          avatar: "",
          status: "offline"
        }
      };
    },
    onSuccess: (session) => {
      setSession(session);
    }
  });

  const errorMessage = useMemo(() => {
    if (!authMutation.error) {
      return "";
    }
    if (authMutation.error instanceof ApiError) {
      return authMutation.error.message;
    }
    return "请求失败，请检查后端服务是否已启动。";
  }, [authMutation.error]);

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    authMutation.mutate();
  };

  return (
    <main className="auth-page">
      <section className="auth-hero" aria-labelledby="auth-title">
        <div className="brand-lockup">
          <div className="brand-mark" aria-hidden="true">
            <MessageCircle size={28} strokeWidth={1.8} />
          </div>
          <span>实时消息工作台</span>
        </div>
        <h1 id="auth-title">Go IM Core</h1>
        <p>
          一个轻量、清爽的聊天前端，用现有 Go 后端完成登录、会话和实时消息。
        </p>
      </section>

      <section className="auth-panel" aria-label="认证表单">
        <div className="auth-tabs" role="group" aria-label="认证方式">
          <button
            className={mode === "login" ? "is-active" : ""}
            type="button"
            onClick={() => setMode("login")}
          >
            登录
          </button>
          <button
            className={mode === "register" ? "is-active" : ""}
            type="button"
            onClick={() => setMode("register")}
          >
            注册
          </button>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          <label>
            用户名
            <input
              autoComplete="username"
              minLength={3}
              required
              value={username}
              onChange={(event) => setUsername(event.target.value)}
            />
          </label>
          {mode === "register" ? (
            <label>
              昵称
              <input
                maxLength={20}
                value={nickname}
                onChange={(event) => setNickname(event.target.value)}
              />
            </label>
          ) : null}
          <label>
            密码
            <input
              autoComplete={mode === "login" ? "current-password" : "new-password"}
              minLength={8}
              required
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
            />
          </label>
          {errorMessage ? <p className="form-error">{errorMessage}</p> : null}
          <button className="button button-primary" disabled={authMutation.isPending}>
            {authMutation.isPending
              ? "提交中"
              : mode === "login"
                ? "进入聊天"
                : "创建并进入"}
          </button>
        </form>
      </section>
    </main>
  );
}
