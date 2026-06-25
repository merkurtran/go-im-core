import { Route, Routes } from "react-router-dom";

import { AuthPage } from "../features/auth/AuthPage";
import { useAuthStore } from "../features/auth/authStore";
import { ChatPage } from "../features/chat/ChatPage";

function RootRoute() {
  const token = useAuthStore((state) => state.token);
  return token ? <ChatPage /> : <AuthPage />;
}

export function AppRouter() {
  return (
    <Routes>
      <Route path="*" element={<RootRoute />} />
    </Routes>
  );
}
