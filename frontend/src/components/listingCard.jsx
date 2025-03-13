// This will be a listing card.
// These cards will be used on the listing page.
// The cards will diplay the item picture, small description
// MAYBE: option to add cart

import Image from "next/image";

const ListingCard = ({ children, imgsrc, ...props }) => {
  // Fallback image URL, e.g., a default image in your public folder
  // Question-Mark
  const fallbackimageSrc = "/Question-Mark.png";

  return (
    <div
      {...props}
      className="group relative h-80 max-w-xs overflow-hidden rounded-2xl shadow-lg"
    >
      {/* The image will a a render a default image a question mark if an image is not proide */}
      <Image
        fill={true}
        src={imgsrc || fallbackimageSrc}
        alt="Listing Image"
        className="h-full w-full object-cover transition-transform duration-200 group-hover:scale-110"
      />
      <div className="absolute inset-0 flex items-end bg-gradient-to-t from-black/60 to-transparent">
        <div className="p-4 text-white">{children}</div>
      </div>
    </div>
  );
};

export default ListingCard;
