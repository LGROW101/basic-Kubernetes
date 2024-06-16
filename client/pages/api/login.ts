import type { NextApiRequest, NextApiResponse } from 'next';
import axios from 'axios';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'POST') {
    const { personalDeduction, kReceipt } = req.body;
    const auth = req.headers.authorization;

    if (!auth) {
      return res.status(401).json({ success: false, message: 'Authorization header is missing' });
    }

    const [username, password] = Buffer.from(auth.split(' ')[1], 'base64').toString().split(':');

    if (username !== process.env.ADMIN_USERNAME || password !== process.env.ADMIN_PASSWORD) {
      return res.status(401).json({ success: false, message: 'Invalid username or password' });
    }

    try {
      const response = await axios.post(
        'http://api:8080/admin/deductions',
        { personalDeduction: personalDeduction || undefined, k_receipt: kReceipt || undefined },
        { headers: { Authorization: auth } }
      );

      const { data } = response;
      const personalDeductionUpdated = personalDeduction !== undefined && data.PersonalDeduction === personalDeduction;
      const kReceiptUpdated = kReceipt !== undefined && data.KReceipt === kReceipt;

      const updatedFields = [];
      if (personalDeductionUpdated) {
        updatedFields.push('Personal Deduction');
      }
      if (kReceiptUpdated) {
        updatedFields.push('K Receipt');
      }

      if (updatedFields.length > 0) {
        const message = `Updated successfully: ${updatedFields.join(', ')}`;
        res.status(200).json({ success: true, message });
      } else {
        res.status(200).json({ success: false, message: 'No fields were updated' });
      }
    } catch (error) {
      console.error('Error updating settings:', error);
      res.status(500).json({ success: false, message: 'Internal server error' });
    }
  } else {
    res.status(405).json({ message: 'Method not allowed' });
  }
}