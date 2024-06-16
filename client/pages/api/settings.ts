// pages/api/settings.ts
import { NextApiRequest, NextApiResponse } from 'next';
import axios, { AxiosError } from 'axios';

interface Setting {
  personalDeduction: number;
  kReceipt: number;
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'GET') {
    try {
      const response = await axios.get<Setting>('http://api:8080/admin/deductions');
      const { personalDeduction, kReceipt } = response.data;
      res.status(200).json({ personalDeduction, kReceipt });
    } catch (error) {
      console.error('Error fetching settings:', error);
      res.status(500).json({ message: 'Internal server error' });
    }
  } else if (req.method === 'POST') {
    const auth = req.headers.authorization;

    if (!auth) {
      return res.status(401).json({ message: 'Authorization header is missing' });
    }

    const [username, password] = Buffer.from(auth.split(' ')[1], 'base64').toString().split(':');

    if (username !== process.env.ADMIN_USERNAME || password !== process.env.ADMIN_PASSWORD) {
      return res.status(401).json({ message: 'Invalid username or password' });
    }

    const { personalDeduction, kReceipt } = req.body;
    try {
      await axios.post<Setting>('http://api:8080/admin/deductions', {
        personalDeduction,
        k_receipt: kReceipt,
      }, {
        headers: { Authorization: auth },
      });
      res.status(200).json({ message: 'Settings updated successfully' });
    } catch (error) {
      console.error('Error updating settings:', error);
      if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError;
        if (axiosError.response && axiosError.response.status === 401) {
          res.status(401).json({ message: 'Unauthorized' });
        } else {
          res.status(500).json({ message: 'Internal server error' });
        }
      } else {
        res.status(500).json({ message: 'Internal server error' });
      }
    }
  } else {
    res.status(405).json({ message: 'Method not allowed' });
  }
}