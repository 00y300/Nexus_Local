// pages/orders/index.jsx
import { useEffect, useState } from "react";
import Link from "next/link";
import { useApi } from "@/context/ApiContext";

export default function OrdersPage() {
  const { getOrders } = useApi();
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    getOrders() // GET http://…/orders
      .then(setOrders)
      .catch(console.error);
  }, []);

  return (
    <div className="p-8">
      <h1 className="mb-4 text-2xl font-bold">All Orders</h1>
      <ul className="space-y-3">
        {orders.map((order) => (
          <li key={order.id} className="rounded border p-4 hover:bg-gray-50">
            <Link
              href={`/orders/${order.id}`}
              className="text-blue-600 hover:underline"
            >
              Order #{order.id} — User {order.user_id} —{" "}
              {new Date(order.created_at).toLocaleString()}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
