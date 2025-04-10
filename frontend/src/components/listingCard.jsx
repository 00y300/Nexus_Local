// components/ListingCard.jsx
import Image from "next/image";
import { useCart } from "@/context/CartContext";

const ListingCard = ({ id, name, description, price, stock, imgsrc }) => {
  const fallbackImage = "/Question‑Mark.png";
  const { addItem } = useCart();
  const isSoldOut = stock <= 0;

  return (
    <div className="group relative h-80 max-w-xs overflow-hidden rounded-2xl shadow-lg">
      <Image
        src={imgsrc || fallbackImage}
        width={500}
        height={500}
        alt={name}
        className="h-full w-full object-cover transition-transform duration-200 group-hover:scale-110"
      />

      <div className="absolute inset-0 flex flex-col justify-end bg-gradient-to-t from-black/60 to-transparent p-4">
        <h2 className="text-lg font-semibold text-white">{name}</h2>
        <p className="text-sm text-gray-200">{description}</p>
        <p className="mt-1 text-xs text-gray-300">
          {isSoldOut ? "Sold Out" : `${stock} in stock`}
        </p>

        <div className="mt-2 flex items-center justify-between">
          <span className="font-bold text-white">${price.toFixed(2)}</span>
          <button
            onClick={() =>
              addItem({ id, name, price, imgsrc: imgsrc || fallbackImage })
            }
            disabled={isSoldOut}
            className={`rounded-md px-3 py-1 text-white ${
              isSoldOut
                ? "cursor-not-allowed bg-gray-500"
                : "bg-blue-600 hover:bg-blue-700"
            }`}
          >
            {isSoldOut ? "—" : "Add to Cart"}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ListingCard;
