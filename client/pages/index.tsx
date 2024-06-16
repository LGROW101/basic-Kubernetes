import { useState, useEffect } from "react";
import axios from "axios";
import router from "next/router";


export default function TaxCalculator() {
  const [totalIncome, setTotalIncome] = useState("");
  const [wht, setWht] = useState("");
  const [kReceiptAmount, setKReceiptAmount] = useState("");
  const [donationAmount, setDonationAmount] = useState("");
  const [personalDeduction, setPersonalDeduction] = useState("");

  useEffect(() => {
    const fetchSettings = async () => {
      try {
        const response = await axios.get("/api/settings");
        const { personalDeduction, kReceipt } = response.data;
        if (personalDeduction !== undefined) {
          setPersonalDeduction(personalDeduction.toString());
        }
        if (kReceipt !== undefined) {
          setKReceiptAmount(kReceipt.toString());
        }
      } catch (error) {
        console.error("Error fetching settings:", error);
      }
    };

    fetchSettings();
  }, []);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const allowances = [
      {
        allowanceType: "k-receipt",
        amount: Number(kReceiptAmount.replace(/[฿,]/g, "")),
      },
      {
        allowanceType: "donation",
        amount: Number(donationAmount.replace(/[฿,]/g, "")),
      },
    ];

    try {
      const response = await axios.post("/api/calculate", {
        totalIncome: Number(totalIncome.replace(/[฿,]/g, "")),
        wht: Number(wht.replace(/[฿,]/g, "")),
        allowances,
        personalDeduction,
        IncludeTaxLevel: true,
      });
      const { level, tax } = response.data;
      router.push(`/TaxResultPage?level=${level}&tax=${tax}`);
    } catch (error) {
      console.error("Error calculating tax:", error);
      alert("เกิดข้อผิดพลาดในการคำนวณภาษี");
    }
  };

  const formatNumber = (value: string) => {
    const formattedValue = value.replace(/[฿,]/g, "");
    return formattedValue !== ""
      ? `฿${Number(formattedValue).toLocaleString()}`
      : "฿";
  };

  return (
    <div className="mb-4  max-w-lg mx-auto p-8">
      <h1 className=" text-center text-2xl font-bold mb-6">
        คำนวณภาษีเงินได้บุคคลธรรมดา
      </h1>
      <form
        onSubmit={handleSubmit}
        className="max-w-lg mx-auto bg-white p-8 rounded-md shadow-md"
      >
        <div className="mb-3">
          <label
            htmlFor="totalIncome"
            className="block text-gray-700 font-semibold mb-2"
          >
            เงินเดือน / ต่อปี
          </label>
          <input
            type="text"
            id="totalIncome"
            value={formatNumber(totalIncome)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setTotalIncome(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <div className="mb-3">
          <label
            htmlFor="wht"
            className="block text-gray-700 font-semibold mb-2"
          >
            ภาษีหัก ณ ที่จ่าย
          </label>
          <input
            type="text"
            id="wht"
            value={formatNumber(wht)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setWht(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <div className="mb-3">
          <label
            htmlFor="kReceiptAmount"
            className="block text-gray-700 font-semibold mb-2"
          >
            โครงการช้อปลดภาษี(k-Receipt)
          </label>
          <input
            type="text"
            id="kReceiptAmount"
            value={formatNumber(kReceiptAmount)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setKReceiptAmount(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <div className="mb-3">
          <label
            htmlFor="donationAmount"
            className="block text-gray-700 font-semibold mb-2"
          >
            เงินบริจาค
          </label>
          <input
            type="text"
            id="donationAmount"
            value={formatNumber(donationAmount)}
            onChange={(e) => {
              const inputValue = e.target.value.replace(/[^0-9]/g, "");
              setDonationAmount(inputValue);
            }}
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
            required
          />
        </div>
        <div className="mb-3">
          <label
            htmlFor="personalDeduction"
            className="block text-gray-700 font-semibold mb-2"
          >
            ค่าลดหย่อนส่วนบุคคล
          </label>
          <input
            type="text"
            disabled
            id="personalDeduction"
            value={
              personalDeduction !== undefined
                ? `฿${personalDeduction.toLocaleString()}`
                : ""
            }
            className="w-full px-3 py-2 text-gray-700 border rounded-md focus:outline-none"
          />
        </div>
        <button
          type="submit"
          className="mt-5 w-full bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600 transition duration-200"
        >
          คำนวณภาษี
        </button>
      </form>
    </div>
  );
}
