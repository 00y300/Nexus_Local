// pages/_app.js
import "@/app/globals.css";
import { CartProvider } from "@/context/CartContext";
import { ApiProvider } from "@/context/ApiContext";
import NavigationBar from "@/components/navigationBar";

export default function MyApp({ Component, pageProps }) {
  return (
    <CartProvider>
      <ApiProvider>
        <NavigationBar />
        <Component {...pageProps} />
      </ApiProvider>

      <footer className="mt-16 border-t py-6 text-center text-sm text-gray-500">
        © {new Date().getFullYear()} Nexus Local. All rights reserved.
      </footer>
    </CartProvider>
  );
}
