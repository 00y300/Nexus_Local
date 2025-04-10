// components/LogoutButton.jsx
import { useState } from "react";
import { useRouter } from "next/router";

const LogoutButton = ({ onLogout }) => {
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleLogout = async () => {
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/logout", {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        // 1) clear parent state
        onLogout();
        // 2) redirect
        router.push("/");
      } else {
        console.error("Logout failed:", res.status);
      }
    } catch (err) {
      console.error("Logout error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <button
      onClick={handleLogout}
      disabled={loading}
      className="rounded bg-red-100 px-4 py-2 text-red-700 hover:bg-red-200 disabled:opacity-50"
    >
      {loading ? "Logging out…" : "Logout"}
    </button>
  );
};

export default LogoutButton;
