import { useState } from "react";
import { useApi } from "@/context/ApiContext";

export default function UpdateItemPage() {
  const { updateItemApi } = useApi();
  const [form, setForm] = useState({
    item_id: "",
    price: "",
    stock: "",
  });
  const [message, setMessage] = useState("");

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        item_id: parseInt(form.item_id, 10),
        ...(form.price !== "" && { price: parseFloat(form.price) }),
        ...(form.stock !== "" && { stock: parseInt(form.stock, 10) }),
      };
      await updateItemApi(payload);
      setMessage(`✅ Updated item ${payload.item_id}`);
      setForm({ item_id: "", price: "", stock: "" });
    } catch (err) {
      setMessage(`❌ ${err.message}`);
    }
  };

  return (
    <div className="mx-auto max-w-md space-y-4 p-8">
      <h1 className="text-2xl font-bold">Update Item</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          name="item_id"
          type="number"
          value={form.item_id}
          onChange={handleChange}
          placeholder="Item ID"
          required
          className="w-full rounded border p-2"
        />
        <input
          name="price"
          type="number"
          step="0.01"
          value={form.price}
          onChange={handleChange}
          placeholder="New Price"
          className="w-full rounded border p-2"
        />
        <input
          name="stock"
          type="number"
          value={form.stock}
          onChange={handleChange}
          placeholder="New Stock"
          className="w-full rounded border p-2"
        />
        <button
          type="submit"
          className="rounded bg-green-600 px-4 py-2 text-white"
        >
          Update Item
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}
