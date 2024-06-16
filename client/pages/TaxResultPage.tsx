

import { useRouter } from 'next/router';

const TaxResultPage = () => {
  const router = useRouter();
  const { tax } = router.query;

  const formatNumber = (value: string) => {
    return Number(value).toLocaleString();
  };

  const handleCalculateAgain = () => {
    router.push('/');
  };

  return (
    <div>
      <div className="mb-4 text-center max-w-lg mx-auto p-8">
        <h1 className="mt-4 w-full px-4 py-2 text-2xl font-bold mb-6">ผลการคำนวณภาษี</h1>
      <div className="max-w-lg mx-auto bg-white p-8 rounded-md shadow-md">
  {tax && (
    <div className=" text-1xl flex justify-between items-center">
      <span>ภาษีที่ต้องชำระ</span>
      <span>{Array.isArray(tax) ? tax.join(', ') : formatNumber(tax)} บาท</span>
    </div>
  )}
</div>
        <button
          type="button"
          onClick={handleCalculateAgain}
          className="mt-5 w-full bg-green-500 text-white px-4 py-2 rounded-md hover:bg-green-600 transition duration-200"
        >
          คำนวณใหม่
        </button>
      </div>
    </div>
  );
};

export default TaxResultPage;