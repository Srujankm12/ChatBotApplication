import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AI Chat",
  description: "AI-powered chat with Gemini 1.5 Flash",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="antialiased">{children}</body>
    </html>
  );
}
