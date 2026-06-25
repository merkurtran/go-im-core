import { apiRequest } from "./client";
import type { UserProfile } from "./types";

export type LoginRequest = {
  username: string;
  password: string;
};

export type RegisterRequest = {
  username: string;
  password: string;
  nickname?: string;
};

export type LoginResponse = {
  token: string;
  user: UserProfile;
};

export type RegisterResponse = {
  user_id: string;
  token: string;
};

export function login(request: LoginRequest) {
  return apiRequest<LoginResponse>("/auth/login", {
    method: "POST",
    body: request
  });
}

export function register(request: RegisterRequest) {
  return apiRequest<RegisterResponse>("/auth/register", {
    method: "POST",
    body: request
  });
}
