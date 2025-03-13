import NavigationBar from "@/components/navigationBar";

export default function MyApp({ Component, pageProps }) {
  return (
    <>
      <NavigationBar></NavigationBar>
      <Component {...pageProps} /> {/* Render the specific page */}
    </>
  );
}
