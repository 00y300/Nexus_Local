// components/NavigationBar.jsx
import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import LoginPopup from "./LoginPopup"; // make sure this file has `export default`
import LogoutButton from "./logoutButton"; // or rename to LogoutButton.jsx
import { useCart } from "@/context/CartContext";
import { motion } from "framer-motion";

export default function NavigationBar() {
  const [showPopup, setShowPopup] = useState(false);
  const [user, setUser] = useState(null);
  const { getTotalItems } = useCart();
  const count = getTotalItems();

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/me`, { credentials: "include" })
      .then((res) => (res.ok ? res.json() : Promise.reject()))
      .then(setUser)
      .catch(() => setUser(null));
  }, []);

  return (
    <motion.nav
      initial={{ y: -50, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.4 }}
      className="sticky top-0 z-50 bg-white/80 shadow-sm backdrop-blur"
    >
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
        {/* Logo / Home link */}
        <Link
          href="/"
          className="text-primary-600 text-2xl font-bold transition hover:opacity-90"
        >
          Nexus Local
        </Link>

        {/* Nav links */}
        <ul className="flex items-center space-x-6">
          {["Home", "Listing", "Orders", "Cart"].map((label) => {
            const href = label === "Home" ? "/" : `/${label.toLowerCase()}`;
            return (
              <li key={label}>
                <Link
                  href={href}
                  className="hover:text-primary-600 relative text-gray-700 transition"
                >
                  {label}
                  {label === "Cart" && count > 0 && (
                    <span className="bg-accent absolute -top-2 -right-3 flex h-5 w-5 items-center justify-center rounded-full text-xs font-semibold text-white">
                      {count}
                    </span>
                  )}
                </Link>
              </li>
            );
          })}
        </ul>

        {/* Login / User */}
        <div className="flex items-center space-x-4">
          {user ? (
            <>
              <span className="text-gray-800">{user.displayName}</span>
              <LogoutButton onLogout={() => setUser(null)} />
            </>
          ) : (
            <button
              onClick={() => setShowPopup(true)}
              className="rounded-full p-1 transition hover:bg-gray-100"
            >
              <Image
                src="/Question-Mark.png"
                width={32}
                height={32}
                alt="Login"
                className="rounded-full border"
              />
            </button>
          )}
          {showPopup && <LoginPopup onClose={() => setShowPopup(false)} />}
        </div>
      </div>
    </motion.nav>
  );
}
