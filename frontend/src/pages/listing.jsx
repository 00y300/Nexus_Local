import ListingCard from "@/components/listingCard";
import { useCart } from "@/context/CartContext";

export default function ListingPage({ items }) {
  const { getTotalItems } = useCart();
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  return (
    <div className="space-y-6 px-4 py-6 sm:px-6 lg:px-8">
      {/* Header */}
      <div className="flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
        <h1 className="text-2xl font-bold sm:text-3xl">Shop Our Products</h1>
        <div className="text-lg">
          Cart:{" "}
          <span className="font-semibold">
            {getTotalItems()} item{getTotalItems() !== 1 && "s"}
          </span>
        </div>
      </div>

      {/* Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
        {items.map((item) => {
          const src = item.image_url
            ? `${apiUrl}${item.image_url}`
            : "/Question-Mark.png";

          return (
            <ListingCard
              key={item.id}
              id={item.id}
              name={item.name}
              description={item.description}
              price={item.price}
              stock={item.stock}
              imgsrc={src}
            />
          );
        })}
      </div>
    </div>
  );
}

export async function getServerSideProps() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items`);
  const items = await res.json();
  return { props: { items } };
}
