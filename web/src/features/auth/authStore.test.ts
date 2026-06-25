import { beforeEach, describe, expect, it } from "vitest";

import { useAuthStore } from "./authStore";

const user = {
  user_id: "u1",
  username: "alice",
  nickname: "Alice",
  avatar: "",
  status: "offline"
};

describe("useAuthStore", () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.setState({ token: null, user: null });
  });

  it("stores and clears an authenticated session", () => {
    useAuthStore.getState().setSession({ token: "jwt", user });

    expect(useAuthStore.getState().token).toBe("jwt");
    expect(useAuthStore.getState().user?.username).toBe("alice");

    useAuthStore.getState().clearSession();

    expect(useAuthStore.getState().token).toBeNull();
    expect(useAuthStore.getState().user).toBeNull();
    expect(localStorage.getItem("go-im-core.session")).toBeNull();
  });

  it("hydrates a saved session from localStorage", () => {
    localStorage.setItem(
      "go-im-core.session",
      JSON.stringify({ token: "saved-token", user })
    );

    useAuthStore.getState().hydrateSession();

    expect(useAuthStore.getState().token).toBe("saved-token");
    expect(useAuthStore.getState().user?.user_id).toBe("u1");
  });
});
