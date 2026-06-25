import type { ProfileUpdate } from "../features/profile/profileTypes";
import { apiRequest } from "./client";
import type { UserProfile } from "./types";

export function getProfile(token: string) {
  return apiRequest<UserProfile>("/users/me", { token });
}

export function updateProfile(token: string, request: ProfileUpdate) {
  return apiRequest<UserProfile | null>("/users/me", {
    method: "PUT",
    token,
    body: request
  });
}
