import { create } from "zustand";

import type { UserProfile } from "../../api/types";

const SESSION_STORAGE_KEY = "go-im-core.session";

type AuthSession = {
  token: string;
  user: UserProfile;
};

type AuthState = {
  token: string | null;
  user: UserProfile | null;
  setSession: (session: AuthSession) => void;
  setUser: (user: UserProfile) => void;
  clearSession: () => void;
  hydrateSession: () => void;
};

export const useAuthStore = create<AuthState>((set, get) => ({
  ...readStoredSession(),
  setSession: (session) => {
    writeStoredSession(session);
    set(session);
  },
  setUser: (user) => {
    const token = get().token;
    if (token) {
      writeStoredSession({ token, user });
    }
    set({ user });
  },
  clearSession: () => {
    removeStoredSession();
    set({ token: null, user: null });
  },
  hydrateSession: () => {
    set(readStoredSession());
  }
}));

function readStoredSession(): Pick<AuthState, "token" | "user"> {
  if (typeof localStorage === "undefined") {
    return { token: null, user: null };
  }

  const raw = localStorage.getItem(SESSION_STORAGE_KEY);
  if (!raw) {
    return { token: null, user: null };
  }

  try {
    const session = JSON.parse(raw) as AuthSession;
    if (!session.token || !session.user?.user_id) {
      return { token: null, user: null };
    }
    return session;
  } catch {
    return { token: null, user: null };
  }
}

function writeStoredSession(session: AuthSession) {
  if (typeof localStorage !== "undefined") {
    localStorage.setItem(SESSION_STORAGE_KEY, JSON.stringify(session));
  }
}

function removeStoredSession() {
  if (typeof localStorage !== "undefined") {
    localStorage.removeItem(SESSION_STORAGE_KEY);
  }
}
