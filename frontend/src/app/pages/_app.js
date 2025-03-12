// This will responsible for adding a default layout for the web app.
import "../styles/index.css";
import NavBar from "../components/navBar"; // Import the NavBar component

export default function MyApp({ Component, pageProps }) {
  return (
    <>
      <NavBar /> {/* Include the NavBar */}
      <Component {...pageProps} /> {/* Render the specific page */}
    </>
  );
}
