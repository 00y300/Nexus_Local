// This is the navigation bar for the Nexus Local
// Utilizes the routes form NEXTJS to link other webpages
// See: https://nextjs.org/docs/app/building-your-application/routing/linking-and-navigating#link-component

// components/NavigationBar.jsx
import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import LoginPopup from "./LoginPopup";
import LogoutButton from "./logoutButton";
import { useCart } from "@/context/CartContext";

const NavigationBar = () => {
  const [showPopup, setShowPopup] = useState(false);
  const [user, setUser] = useState(null);
  const { getTotalItems } = useCart();
  const count = getTotalItems();

  // On mount, try to fetch the current user
  useEffect(() => {
    fetch("http://localhost:8080/me", {
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("not logged in");
        return res.json();
      })
      .then((data) => setUser(data))
      .catch(() => setUser(null));
  }, []);

  return (
    <nav className="flex h-16 items-center justify-between bg-gray-100 px-4">
      {/* Left side links */}
      <ul className="flex space-x-4">
        <li className="border-2 p-2">
          <Link href="/">Home</Link>
        </li>
        <li className="border-2 p-2">
          <Link href="/listing">Listing</Link>
        </li>
        <li className="border-2 p-2">
          <Link href="/orders">Orders</Link>
        </li>
        <li className="relative border-2 p-2">
          <Link href="/cart">Cart</Link>
          {count > 0 && (
            <span className="absolute -top-2 -right-2 inline-flex h-5 w-5 items-center justify-center rounded-full bg-red-600 text-xs text-white">
              {count}
            </span>
          )}
        </li>
      </ul>

      {/* Right side: login / user display + logout */}
      <div className="relative flex items-center space-x-2">
        {user ? (
          <>
            <span className="px-4 py-2 font-medium text-gray-800">
              {user.displayName}
            </span>
            <LogoutButton onLogout={() => setUser(null)} />
          </>
        ) : (
          <button onClick={() => setShowPopup(true)}>
            <Image
              width={32}
              height={32}
              src="/Question-Mark.png"
              alt="Login"
              className="rounded-full"
            />
          </button>
        )}

        {showPopup && <LoginPopup onClose={() => setShowPopup(false)} />}
      </div>
    </nav>
  );
};

export default NavigationBar;
