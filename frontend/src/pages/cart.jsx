// pages/cart.jsx
import Link from "next/link";
import Image from "next/image";
import { useCart } from "@/context/CartContext";

export default function CartPage() {
  const {
    cartItems,
    increaseQuantity,
    decreaseQuantity,
    removeItem,
    clearCart,
    getTotalPrice,
  } = useCart();

  // Empty state
  if (cartItems.length === 0) {
    return (
      <div className="p-8 text-center">
        <h1 className="mb-4 text-2xl font-bold">Your Cart is Empty</h1>
        <Link href="/listing">
          <a className="text-blue-600 hover:underline">Continue Shopping</a>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6 p-8">
      <h1 className="text-2xl font-bold">Your Cart</h1>

      <ul className="space-y-4">
        {cartItems.map((item) => (
          <li
            key={item.id}
            className="flex items-center justify-between border-b pb-4"
          >
            <div className="flex items-center space-x-4">
              <div className="relative h-20 w-20">
                <Image
                  src={item.imgsrc}
                  alt={item.name}
                  layout="fill"
                  objectFit="cover"
                  className="rounded-lg"
                />
              </div>
              <div>
                <h2 className="text-lg font-semibold">{item.name}</h2>
                <p className="text-sm text-gray-600">
                  ${item.price.toFixed(2)}
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-2">
              <button
                onClick={() => decreaseQuantity(item.id)}
                className="rounded bg-gray-200 px-2 py-1"
              >
                −
              </button>
              <span>{item.quantity}</span>
              <button
                onClick={() => increaseQuantity(item.id)}
                className="rounded bg-gray-200 px-2 py-1"
              >
                +
              </button>
              <button
                onClick={() => removeItem(item.id)}
                className="ml-4 text-red-600 hover:underline"
              >
                Remove
              </button>
            </div>
          </li>
        ))}
      </ul>

      <div className="flex items-center justify-between pt-4">
        <button onClick={clearCart} className="text-red-600 hover:underline">
          Clear Cart
        </button>
        <div className="text-xl font-bold">
          Total: ${getTotalPrice().toFixed(2)}
        </div>
      </div>
    </div>
  );
}
