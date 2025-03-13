import Link from "next/link";

const NavigationBar = () => {
  return (
    <div className="flex h-16 items-center justify-center">
      <ul className="flex space-x-4">
        <li className="border-2 p-2">
          <Link href="/">Home</Link>
        </li>
      </ul>
    </div>
  );
};
export default NavigationBar;
