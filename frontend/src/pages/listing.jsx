// pages/listing.jsx
import ListingCard from "@/components/listingCard";
import { useCart } from "@/context/CartContext";

export default function ListingPage({ items }) {
  const { getTotalItems } = useCart();

  return (
    <div className="space-y-6 p-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Shop Our Products</h1>
        <div>
          Cart:{" "}
          <span className="font-semibold">
            {getTotalItems()} item{getTotalItems() !== 1 && "s"}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 justify-items-center gap-6 sm:grid-cols-2 md:grid-cols-3">
        {items.map((item) => (
          <ListingCard
            key={item.id}
            id={item.id}
            name={item.name}
            description={item.description}
            price={item.price}
            stock={item.stock}
            // uncomment when you have an image URL
            // imgsrc={item.image_url}
          />
        ))}
      </div>
    </div>
  );
}

export async function getServerSideProps() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/items`);
  const items = await res.json();
  return { props: { items } };
}
