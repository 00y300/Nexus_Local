// components/ListingCard.jsx
import Image from "next/image";
import { useCart } from "@/context/CartContext";

export default function ListingCard({
  id,
  name,
  description,
  price,
  stock,
  imgsrc,
}) {
  const fallbackImage = "/Question‑Mark.png";
  const { addItem } = useCart();
  const isSoldOut = stock <= 0;

  return (
    <div className="group flex flex-col overflow-hidden rounded-2xl bg-white shadow-lg">
      {/* Image wrapper: maintains a square aspect ratio */}
      <div className="relative h-0 w-full overflow-hidden pb-[100%]">
        <Image
          src={imgsrc || fallbackImage}
          alt={name}
          layout="fill"
          objectFit="cover"
          className="transform transition-transform duration-200 group-hover:scale-110"
        />
      </div>

      {/* Content */}
      <div className="flex flex-1 flex-col justify-between p-4">
        <div>
          <h2 className="truncate text-lg font-semibold text-gray-900">
            {name}
          </h2>
          <p className="mt-1 line-clamp-2 text-sm text-gray-600">
            {description}
          </p>
          <p className="mt-2 text-xs text-gray-500">
            {isSoldOut ? "Sold Out" : `${stock} in stock`}
          </p>
        </div>

        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-bold text-gray-900">
            ${price.toFixed(2)}
          </span>
          <button
            onClick={() =>
              addItem({ id, name, price, imgsrc: imgsrc || fallbackImage })
            }
            disabled={isSoldOut}
            className={`rounded-md px-3 py-1 text-sm font-medium text-white transition ${
              isSoldOut
                ? "cursor-not-allowed bg-gray-400"
                : "bg-blue-600 hover:bg-blue-700"
            }`}
          >
            {isSoldOut ? "—" : "Add to Cart"}
          </button>
        </div>
      </div>
    </div>
  );
}
