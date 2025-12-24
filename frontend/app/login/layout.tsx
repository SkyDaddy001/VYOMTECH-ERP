import type { Metadata } from 'next';
import { Geist } from 'next/font/google';
import '../globals.css';

const geist = Geist({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Login - VYOM ERP',
  description: 'Sign in to your VYOM ERP account',
};

export default function LoginLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={geist.className}>
        {children}
      </body>
    </html>
  );
}
