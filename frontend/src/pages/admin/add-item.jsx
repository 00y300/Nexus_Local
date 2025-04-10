import { useState } from "react";
import { useApi } from "@/context/ApiContext";

export default function AddItemPage() {
  const { addItemApi } = useApi();
  const [form, setForm] = useState({
    name: "",
    description: "",
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
        name: form.name,
        description: form.description,
        price: parseFloat(form.price),
        stock: parseInt(form.stock, 10),
      };
      const data = await addItemApi(payload);
      setMessage(`✅ Added item with ID ${data.id}`);
      setForm({ name: "", description: "", price: "", stock: "" });
    } catch (err) {
      setMessage(`❌ ${err.message}`);
    }
  };

  return (
    <div className="mx-auto max-w-md space-y-4 p-8">
      <h1 className="text-2xl font-bold">Add New Item</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          name="name"
          value={form.name}
          onChange={handleChange}
          placeholder="Name"
          required
          className="w-full rounded border p-2"
        />
        <textarea
          name="description"
          value={form.description}
          onChange={handleChange}
          placeholder="Description"
          required
          className="w-full rounded border p-2"
        />
        <input
          name="price"
          type="number"
          step="0.01"
          value={form.price}
          onChange={handleChange}
          placeholder="Price"
          required
          className="w-full rounded border p-2"
        />
        <input
          name="stock"
          type="number"
          value={form.stock}
          onChange={handleChange}
          placeholder="Stock"
          required
          className="w-full rounded border p-2"
        />
        <button
          type="submit"
          className="rounded bg-blue-600 px-4 py-2 text-white"
        >
          Add Item
        </button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}
