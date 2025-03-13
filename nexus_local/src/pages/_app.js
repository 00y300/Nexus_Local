// This will be the layout page
// Essentially any components that will be render on every page
import "../app/globals.css";
import NavigationBar from "@/components/navigationBar";

export default function MyApp({ Component, pageProps }) {
  return (
    <>
      <NavigationBar></NavigationBar>
      <Component {...pageProps} /> {/* Render the specific page */}
    </>
  );
}
