// This is the navigation bar for the Nexus Local
// Utilizes the routes form NEXTJS to link other webpages
// See: https://nextjs.org/docs/app/building-your-application/routing/linking-and-navigating#link-component

import Link from "next/link";
import { useState } from "react";
import LoginPopup from "./LoginPopup"; // Adjust the import path as necessary
import Image from "next/image";

const NavigationBar = () => {
  const [showPopup, setShowPopup] = useState(false);

  const handleProfileClick = () => {
    setShowPopup(true);
  };

  const handleClosePopup = () => {
    setShowPopup(false);
  };

  return (
    <nav className="flex h-16 items-center justify-between bg-gray-100 px-4">
      {/* Navigation links */}
      <ul className="flex space-x-4">
        <li className="border-2 p-2">
          <Link href="/">Home</Link>
        </li>
        <li className="border-2 p-2">
          <Link href="/listing">Listing</Link>
        </li>
      </ul>

      {/* Profile icon */}
      <div className="relative">
        <button onClick={handleProfileClick}>
          <Image
            width={500}
            height={500}
            src="/Question-Mark.png"
            alt="Profile"
            className="h-8 w-8 rounded-full"
          />
        </button>

        {/* Login Popup */}
        {showPopup && <LoginPopup onClose={handleClosePopup} />}
      </div>
    </nav>
  );
};

export default NavigationBar;
