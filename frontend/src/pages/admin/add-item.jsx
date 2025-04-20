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
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState("");

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const handleFileChange = (e) => {
    if (e.target.files.length > 0) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const data = new FormData();
    data.append("name", form.name);
    data.append("description", form.description);
    data.append("price", form.price);
    data.append("stock", form.stock);
    if (file) data.append("image", file); // key matches Go handler’s FormFile("image")

    try {
      const json = await addItemApi(data);
      setMessage(`✅ Added item #${json.item_id}`);
      setForm({ name: "", description: "", price: "", stock: "" });
      setFile(null);
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
        <input
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="w-full rounded border p-2"
        />

        <button
          type="submit"
          className="w-full rounded bg-blue-600 px-4 py-2 text-white"
        >
          Add Item with Image
        </button>
      </form>
      {message && <p className="mt-2">{message}</p>}
    </div>
  );
}
