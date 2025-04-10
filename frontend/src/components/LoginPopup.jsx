// components/LoginPopup.jsx
import Image from "next/image";

const LoginPopup = ({ onClose }) => (
  <div
    id="login-popup"
    tabIndex="-1"
    className="fixed top-0 right-0 left-0 z-50 flex h-full items-center justify-center overflow-x-hidden overflow-y-auto bg-black/50"
  >
    <div className="relative h-full w-full max-w-md p-4 md:h-auto">
      <div className="relative rounded-lg bg-white shadow">
        <button
          type="button"
          onClick={onClose}
          className="absolute top-3 right-2.5 ml-auto inline-flex items-center rounded-lg bg-transparent p-1.5 text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900"
        >
          <svg
            aria-hidden="true"
            className="h-5 w-5"
            fill="#c6c7c7"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fillRule="evenodd"
              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0
                 111.414 1.414L11.414 10l4.293 4.293a1 1 0
                 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0
                 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0
                 010-1.414z"
              clipRule="evenodd"
            />
          </svg>
          <span className="sr-only">Close popup</span>
        </button>

        <div className="p-5">
          <div className="text-center">
            <p className="mb-3 text-2xl leading-5 font-semibold text-slate-900">
              Login to your account
            </p>
          </div>

          <div className="mt-7 flex flex-col gap-2">
            <button
              type="button"
              onClick={() => {
                // build a redirect back to wherever the user is right now
                const redirectUri = encodeURIComponent(window.location.href);
                window.location.assign(
                  `http://localhost:8080/login?redirect_uri=${redirectUri}`,
                );
              }}
              className="inline-flex h-10 w-full items-center justify-center gap-2 rounded border border-slate-300 bg-white p-2 text-sm font-medium text-black focus:ring-2 focus:ring-[#333] focus:ring-offset-1"
            >
              <Image
                src="/Microsoft_logo.svg"
                alt="Microsoft"
                width={18}
                height={18}
                className="h-[18px] w-[18px]"
              />
              Continue with Microsoft
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
);

export default LoginPopup;
