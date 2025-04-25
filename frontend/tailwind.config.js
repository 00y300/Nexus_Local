// tailwind.config.js
import defaultTheme from "tailwindcss/defaultTheme";

export default {
  content: ["./src/**/*.{js,jsx,ts,tsx}", "./pages/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", ...defaultTheme.fontFamily.sans],
      },
      colors: {
        primary: {
          50: "#e3f2ff",
          100: "#b3daff",
          200: "#81c2ff",
          300: "#4faaff",
          400: "#1d92ff",
          500: "#006fd6",
          600: "#0054a0",
          700: "#00396b",
          800: "#001f35",
          900: "#00030a",
        },
        accent: {
          DEFAULT: "#FF6B6B",
          light: "#FF8E8E",
          dark: "#E64545",
        },
      },
      borderRadius: {
        xl: "1rem",
        "2xl": "1.5rem",
      },
      boxShadow: {
        card: "0 4px 24px rgba(0, 0, 0, 0.05)",
      },
    },
  },
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
};
