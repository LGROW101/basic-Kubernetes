import { NextApiRequest, NextApiResponse } from "next";

interface CalculationData {
  // ประกาศโครงสร้างข้อมูลของ data ที่จะส่งไปยัง API
  totalIncome: number;
  wht: number;
  allowances: {
    allowanceType: string;
    amount: number;
  }[];
  personalDeduction: number;
}

interface CalculationResponse {
  tax: number;
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'POST') {
    const data: CalculationData = req.body;

    try {
      const tax = await calculateTax(data);
      res.status(200).json({ tax });
    } catch (error) {
      console.error('Error calculating tax:', error);
      res.status(500).json({ error: 'Failed to calculate tax' });
    }
  } else {
    res.status(405).json({ error: 'Method not allowed' });
  }
}

async function calculateTax(data: CalculationData): Promise<number> {
  const response = await fetch('http://api:8080/tax/calculations', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error('Failed to calculate tax');
  }

  const responseData: CalculationResponse = await response.json();
  return responseData.tax;
}