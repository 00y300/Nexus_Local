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
    </CartProvider>
  );
}
