// pages/api/auth.ts
import type { NextApiRequest, NextApiResponse } from 'next';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === 'GET') {
    const auth = req.headers.authorization;

    if (!auth) {
      return res.status(401).json({ message: 'Authorization header is missing' });
    }

    const [username, password] = Buffer.from(auth.split(' ')[1], 'base64').toString().split(':');

    if (username === process.env.ADMIN_USERNAME && password === process.env.ADMIN_PASSWORD) {
      res.status(200).end();
    } else {
      res.status(401).json({ message: 'Invalid username or password' });
    }
  } else {
    res.status(405).end();
  }
}