import type { Metadata } from "next";
import { Geist, Geist_Mono, Halant } from "next/font/google";
import "./globals.css";
import Providers from "./providers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const halant = Halant({ 
    weight: ['300', '400', '500', '600', '700'],
    subsets: ['latin'],
    variable: '--font-halant'
})

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Upwell Grant Manager",
  description: "Upwell's Grant management software written for grant writers to manage grants they've researched, and to collaborate with their clients on grant applications.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${halant.variable} antialiased`}
      >
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
