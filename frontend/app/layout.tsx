import type { Metadata } from 'next';
import { Geist } from 'next/font/google';
import './globals.css';
import Navigation from '@/components/Navigation';
import { ProtectedLayout } from '@/components/ProtectedLayout';

const geist = Geist({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'VYOM ERP - User Count & Seat Management',
  description: 'Enterprise Resource Planning System with User Count and Seat Management',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={geist.className}>
        <main>
          <ProtectedLayout>{children}</ProtectedLayout>
        </main>
      </body>
    </html>
  );
}
