import {
  FileText,
  Home,
  LogOut,
  MessageCircle,
  Plus,
  Search,
  Settings,
} from "lucide-react";
import { useAuthStore } from "../../auth/authStore";

export function LeftNavBar() {
  const user = useAuthStore((s) => s.user);
  const clearSession = useAuthStore((s) => s.clearSession);

  const initial = (user?.nickname || user?.username || "U")
    .trim()
    .slice(0, 1)
    .toUpperCase();

  return (
    <nav className="left-nav" aria-label="主导航">
      <div className="nav-top">
        <div className="nav-logo" aria-label="logo">
          <div className="logo-dots">
            <span />
            <span />
            <span />
            <span />
          </div>
        </div>

        <button className="nav-btn" type="button" aria-label="搜索">
          <Search size={20} strokeWidth={1.8} />
        </button>
        <button className="nav-btn" type="button" aria-label="主页">
          <Home size={20} strokeWidth={1.8} />
        </button>
        <button className="nav-btn" type="button" aria-label="新建">
          <Plus size={20} strokeWidth={1.8} />
        </button>
        <button className="nav-btn" type="button" aria-label="文档">
          <FileText size={20} strokeWidth={1.8} />
        </button>
      </div>

      <div className="nav-bottom">
        <button className="nav-btn" type="button" aria-label="消息">
          <MessageCircle size={20} strokeWidth={1.8} />
        </button>
        <button className="nav-btn" type="button" aria-label="设置">
          <Settings size={20} strokeWidth={1.8} />
        </button>
        <button className="nav-btn avatar-btn" type="button" aria-label="个人资料">
          {initial}
        </button>
        <button
          className="nav-btn"
          type="button"
          aria-label="退出登录"
          onClick={clearSession}
        >
          <LogOut size={20} strokeWidth={1.8} />
        </button>
      </div>
    </nav>
  );
}
