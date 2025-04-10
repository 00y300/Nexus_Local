// pages/_app.js
import "@/app/globals.css";
import { CartProvider } from "@/context/CartContext";
import NavigationBar from "@/components/navigationBar";

export default function MyApp({ Component, pageProps }) {
  return (
    <CartProvider>
      <NavigationBar />
      <Component {...pageProps} />
    </CartProvider>
  );
}
