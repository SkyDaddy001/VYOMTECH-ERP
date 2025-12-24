import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'VYOM LMS - Lead Management System',
  description: 'Enterprise Lead Management System for Sales Teams',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
