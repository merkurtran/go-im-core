import { beforeEach, describe, expect, it, vi } from "vitest";

import { apiRequest, parseEnvelope } from "./client";

describe("parseEnvelope", () => {
  it("returns data for successful backend envelopes", async () => {
    await expect(
      parseEnvelope({ code: 0, message: "success", data: { ok: true } })
    ).resolves.toEqual({ ok: true });
  });

  it("throws an ApiError for backend business errors", async () => {
    await expect(
      parseEnvelope({ code: 1001, message: "invalid request", data: null })
    ).rejects.toMatchObject({
      code: 1001,
      message: "invalid request"
    });
  });
});

describe("apiRequest", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("adds a bearer token when provided", async () => {
    const fetchMock = vi.fn(async () => {
      return new Response(
        JSON.stringify({ code: 0, message: "success", data: { user_id: "u1" } }),
        { status: 200, headers: { "Content-Type": "application/json" } }
      );
    });
    vi.stubGlobal("fetch", fetchMock);

    await apiRequest("/users/me", { token: "abc" });

    expect(fetchMock).toHaveBeenCalledWith(
      "http://localhost:8080/api/v1/users/me",
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: "Bearer abc"
        })
      })
    );
  });
});
