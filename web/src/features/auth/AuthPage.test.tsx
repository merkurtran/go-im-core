import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it } from "vitest";

import { AuthPage } from "./AuthPage";

function renderAuthPage() {
  const queryClient = new QueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter
        future={{ v7_relativeSplatPath: true, v7_startTransition: true }}
      >
        <AuthPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe("AuthPage", () => {
  it("renders login and register actions", () => {
    renderAuthPage();

    expect(
      screen.getByRole("heading", { name: /go im core/i })
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /登录/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /注册/i })).toBeInTheDocument();
  });
});
