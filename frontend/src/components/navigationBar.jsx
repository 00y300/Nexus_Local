// This is the navigation bar for the Nexus Local
// Utilizes the routes form NEXTJS to link other webpages
// See: https://nextjs.org/docs/app/building-your-application/routing/linking-and-navigating#link-component

import Link from "next/link";
import { useState } from "react";
import LoginPopup from "./LoginPopup";
import Image from "next/image";
import { useCart } from "@/context/CartContext";

const NavigationBar = () => {
  const [showPopup, setShowPopup] = useState(false);
  const { getTotalItems } = useCart();
  const count = getTotalItems();

  return (
    <nav className="flex h-16 items-center justify-between bg-gray-100 px-4">
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
      <div className="relative">
        <button onClick={() => setShowPopup(true)}>
          <Image
            width={32}
            height={32}
            src="/Question-Mark.png"
            alt="Profile"
            className="rounded-full"
          />
        </button>
        {showPopup && <LoginPopup onClose={() => setShowPopup(false)} />}
      </div>
    </nav>
  );
};

export default NavigationBar;
