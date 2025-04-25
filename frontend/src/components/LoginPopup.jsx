// components/LoginPopup.jsx
import { useState, useEffect } from "react";
import { createPortal } from "react-dom";
import { motion, AnimatePresence } from "framer-motion";
import Image from "next/image";

export default function LoginPopup({ onClose }) {
  // Only render portal on the client
  const [mounted, setMounted] = useState(false);
  useEffect(() => {
    setMounted(true);
  }, []);
  if (!mounted) return null;

  return createPortal(
    <AnimatePresence>
      {/* Backdrop */}
      <motion.div
        key="backdrop"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.2 }}
        className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 px-4 backdrop-blur-sm"
        onClick={onClose}
        aria-modal="true"
        role="dialog"
      >
        {/* Modal panel */}
        <motion.div
          key="modal"
          initial={{ scale: 0.9, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.9, opacity: 0 }}
          transition={{ type: "spring", stiffness: 300, damping: 25 }}
          className="relative w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Close button */}
          <button
            onClick={onClose}
            className="absolute top-4 right-4 inline-flex items-center justify-center rounded-full p-1.5 text-gray-400 transition hover:bg-gray-100 hover:text-gray-600"
            aria-label="Close login modal"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>

          {/* Content */}
          <div className="mt-2 space-y-4">
            <h2 className="text-center text-2xl font-semibold text-gray-900">
              Sign in to Your Account
            </h2>
            <button
              type="button"
              onClick={() => {
                const redirectUri = encodeURIComponent(window.location.href);
                window.location.assign(
                  `http://localhost:8080/login?redirect_uri=${redirectUri}`,
                );
              }}
              className="focus:ring-primary-300 flex w-full items-center justify-center gap-2 rounded-lg border border-gray-200 bg-white py-3 text-sm font-medium text-gray-800 shadow-sm transition hover:bg-gray-50 focus:ring-2 focus:outline-none"
            >
              <Image
                src="/Microsoft_logo.svg"
                alt="Microsoft"
                width={20}
                height={20}
                className="h-5 w-5"
              />
              Continue with Microsoft
            </button>
          </div>
        </motion.div>
      </motion.div>
    </AnimatePresence>,
    document.body,
  );
}
