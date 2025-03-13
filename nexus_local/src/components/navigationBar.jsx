// This is the navigation bar for the Nexus Local
// Utilizes the routes form NEXTJS to link other webpages
// See: https://nextjs.org/docs/app/building-your-application/routing/linking-and-navigating#link-component
import Link from "next/link";

const NavigationBar = () => {
  return (
    <div className="flex h-16 items-center justify-center">
      <ul className="flex space-x-4">
        <li className="border-2 p-2">
          <Link href="/">Home</Link>
        </li>

        <li className="border-2 p-2">
          <Link href="/listing">Listing</Link>
        </li>
      </ul>
    </div>
  );
};
export default NavigationBar;
