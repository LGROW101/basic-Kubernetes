// /pages/admin/login.tsx
import { useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";

export default function AdminLogin() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [personalDeduction, setPersonalDeduction] = useState("");
  const [kReceipt, setKReceipt] = useState("");
  const [message, setMessage] = useState("");
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Username:", username);
    console.log("Password:", password);
    console.log("Personal Deduction:", personalDeduction);
    console.log("K Receipt:", kReceipt);

    try {
      const auth =
        "Basic " + Buffer.from(`${username}:${password}`).toString("base64");

      const response = await axios.post(
        "/api/login",
        {
          personalDeduction: Number(personalDeduction),
          kReceipt: Number(kReceipt),
        },
        { headers: { Authorization: auth } }
      );

      if (response.data.success) {
        setMessage(response.data.message);
        router.push("/");
      } else {
        setMessage(response.data.message);
      }
    } catch (error) {
      console.error("Error during login:", error);
      setMessage("An error occurred during login");
    }
  };
  const formatNumber = (value: string) => {
    if (value === "") {
      return "฿";
    }
    return `฿${Number(value).toLocaleString()}`;
  };
  return (
    <div className="max-w-sm mx-auto mt-8">
      {message && <p className=" text-center text-red-500 mb-4 text-2xl">{message}</p>}
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="username" className="block mb-2 font-bold">
            Username:
          </label>
          <input
            type="text"
            id="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="password" className="block mb-2 font-bold">
            Password:
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="personalDeduction" className="block mb-2 font-bold">
            ค่าลดหย่อนส่วนตัว:
          </label>
          <input
            type="text"
            id="personalDeduction"
            value={formatNumber(personalDeduction)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setPersonalDeduction(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="kReceipt" className="block mb-2 font-bold">
            ค่าเบี้ยเลี้ยง:
          </label>
          <input
            type="text"
            id="kReceipt"
            value={formatNumber(kReceipt)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setKReceipt(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <button
          type="submit"
          className="mt-5 w-full bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600 transition duration-200"
        >
          Login
        </button>
      </form>
    </div>
  );
}
