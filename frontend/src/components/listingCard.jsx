// src/components/listingCard.jsx
import Image from "next/image";
import { useCart } from "@/context/CartContext";
import { motion } from "framer-motion";

export default function ListingCard({
  id,
  name,
  description,
  price,
  stock = 1,
  imgsrc,
}) {
  const { addItem } = useCart();
  const soldOut = stock <= 0;
  const fallback = "/Question-Mark.png";

  return (
    <motion.div
      className="group flex flex-col overflow-hidden rounded-2xl bg-white shadow-lg"
      whileHover={{ scale: 1.02, boxShadow: "0 8px 24px rgba(0,0,0,0.1)" }}
      transition={{ type: "spring", stiffness: 300 }}
    >
      {/* 1. Image */}
      <div className="relative h-0 w-full overflow-hidden pb-[100%]">
        <Image
          src={imgsrc || fallback}
          alt={name}
          fill
          className="object-cover transition-transform duration-200 group-hover:scale-110"
        />
      </div>

      {/* 2. Details */}
      <div className="flex flex-1 flex-col justify-between p-4">
        <div>
          <h3 className="truncate text-lg font-semibold text-gray-900">
            {name}
          </h3>
          <p className="mt-1 line-clamp-2 text-sm text-gray-600">
            {description}
          </p>
          <p className="mt-2 text-xs text-gray-500">
            {soldOut ? "Sold Out" : `${stock} in stock`}
          </p>
        </div>

        {/* 3. Price + Add to Cart */}
        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-bold text-gray-900">
            ${price.toFixed(2)}
          </span>

          <button
            onClick={() =>
              addItem({ id, name, price, imgsrc: imgsrc || fallback })
            }
            disabled={soldOut}
            className={`rounded-lg px-3 py-1 text-sm font-medium transition ${
              soldOut
                ? "cursor-not-allowed bg-gray-200 text-gray-400"
                : "bg-blue-600 text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-300 focus:outline-none"
            }`}
          >
            {soldOut ? "—" : "Add to Cart"}
          </button>
        </div>
      </div>
    </motion.div>
  );
}
