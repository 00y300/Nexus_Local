import { createContext, useContext } from "react";

const ApiContext = createContext();

export const ApiProvider = ({ children }) => {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  // now accepts FormData (for file uploads)
  const addItemApi = async (formData) => {
    const res = await fetch(`${apiUrl}/items/add`, {
      method: "POST",
      credentials: "include", // include cookies
      body: formData, // browser sets multipart boundaries
    });
    if (!res.ok) throw new Error(`Add item failed: ${res.statusText}`);
    return res.json();
  };

  const updateItemApi = async (data) => {
    const res = await fetch(`${apiUrl}/items/update`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(`Update item failed: ${res.statusText}`);
    return res.json();
  };

  const getOrders = async (order_id) => {
    const url = order_id
      ? `${apiUrl}/orders?order_id=${order_id}`
      : `${apiUrl}/orders`;
    const res = await fetch(url, {
      credentials: "include",
    });
    if (!res.ok) throw new Error(`Fetch orders failed: ${res.statusText}`);
    return res.json();
  };

  const postOrderApi = async (order) => {
    const res = await fetch(`${apiUrl}/orders`, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(order),
    });
    if (!res.ok) throw new Error(`Post order failed: ${res.statusText}`);
    return res.json();
  };

  const deleteOrderApi = async (order_id) => {
    const res = await fetch(`${apiUrl}/orders?order_id=${order_id}`, {
      method: "DELETE",
      credentials: "include",
    });
    if (!res.ok) throw new Error(`Delete order failed: ${res.statusText}`);
    return res.text();
  };

  return (
    <ApiContext.Provider
      value={{
        addItemApi,
        updateItemApi,
        getOrders,
        postOrderApi,
        deleteOrderApi,
      }}
    >
      {children}
    </ApiContext.Provider>
  );
};

export const useApi = () => useContext(ApiContext);
