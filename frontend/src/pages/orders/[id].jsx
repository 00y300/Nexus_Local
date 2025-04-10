// pages/orders/[id].jsx
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/router";

export default function OrderDetailPage({ order }) {
  const router = useRouter();

  // In case of client‑side navigation
  if (!order) {
    return <p className="p-8">Loading…</p>;
  }

  const handleDelete = async () => {
    await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/orders?order_id=${order.id}`,
      { method: "DELETE" },
    );
    router.push("/orders");
  };

  return (
    <div className="space-y-4 p-8">
      <h1 className="text-2xl font-bold">Order #{order.id}</h1>
      <p>User: {order.user_id}</p>
      <p>Placed at: {new Date(order.created_at).toLocaleString()}</p>

      <h2 className="mt-4 text-lg font-semibold">Items</h2>
      <ul className="space-y-2">
        {order.items.map(({ id, name, price, quantity }) => (
          <li
            key={id}
            className="flex items-center justify-between border-b pb-2"
          >
            <div className="flex items-center space-x-4">
              <Image
                src={name ? `/images/${id}.jpg` : "/Question‑Mark.png"}
                alt={name}
                width={60}
                height={60}
                className="rounded"
              />
              <div>
                <div className="font-semibold">{name}</div>
                <div className="text-sm text-gray-600">
                  ${price.toFixed(2)} × {quantity} ={" "}
                  <span className="font-bold">
                    ${(price * quantity).toFixed(2)}
                  </span>
                </div>
              </div>
            </div>
          </li>
        ))}
      </ul>

      <div className="flex items-center justify-between pt-4">
        <button onClick={handleDelete} className="text-red-600 hover:underline">
          Delete Order
        </button>
        <Link href="/orders" className="text-blue-600 hover:underline">
          ← Back to Orders
        </Link>
      </div>
    </div>
  );
}

export async function getServerSideProps({ params }) {
  const api = process.env.NEXT_PUBLIC_API_URL;
  const { id } = params;

  // 1) Fetch the order + its items
  const ordRes = await fetch(`${api}/orders?order_id=${id}`);
  if (!ordRes.ok) return { notFound: true };
  const { order: ord, order_items } = await ordRes.json();

  // 2) Fetch full product catalog
  const catRes = await fetch(`${api}/items`);
  const catalog = await catRes.json();

  // 3) Merge each line with product details
  const merged = order_items.map(({ item_id, quantity }) => {
    const prod = catalog.find((p) => p.id === item_id) || {};
    return {
      id: item_id,
      name: prod.name || `#${item_id}`,
      price: prod.price || 0,
      quantity,
    };
  });

  return {
    props: {
      order: {
        id: ord.id,
        user_id: ord.user_id,
        created_at: ord.created_at,
        items: merged,
      },
    },
  };
}
