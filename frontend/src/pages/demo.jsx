// pages/demo.jsx
import ListingCard from "@/components/listingCard";
import { useCart } from "@/context/CartContext";

export default function DemoPage() {
  const { cartItems, addItem, removeItem, getTotalItems, getTotalPrice } =
    useCart();

  // A fake product we’ll add to cart when you click the button
  const demoProduct = {
    id: "demo-1",
    name: "Demo Widget",
    price: 19.99,
    imgsrc: "/widget.png",
    description: "This is just a demo widget.",
  };

  return (
    <div className="space-y-8 p-8">
      <h1 className="text-3xl font-bold">Cart Demo</h1>

      {/* 1) A ListingCard you can “Add to Cart” */}
      <ListingCard {...demoProduct} />

      {/* 2) A custom button doing the same thing programmatically */}
      <button
        onClick={() => addItem(demoProduct)}
        className="mt-4 rounded-md bg-green-600 px-4 py-2 text-white"
      >
        Add Demo Widget Programmatically
      </button>

      {/* 3) Show cart summary */}
      <div className="mt-8 border-t pt-4">
        <h2 className="text-2xl font-semibold">Cart Summary</h2>
        <p>Total Items: {getTotalItems()}</p>
        <p>Total Price: ${getTotalPrice().toFixed(2)}</p>

        {/* 4) List out current items with remove buttons */}
        {cartItems.map((item) => (
          <div
            key={item.id}
            className="mt-2 flex items-center justify-between rounded border p-2"
          >
            <div>
              <strong>{item.name}</strong> x {item.quantity} = $
              {(item.price * item.quantity).toFixed(2)}
            </div>
            <button
              onClick={() => removeItem(item.id)}
              className="text-red-600 hover:underline"
            >
              Remove
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
