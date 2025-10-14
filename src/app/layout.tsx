import type { Metadata } from "next";
import "./globals.css";
import React from "react";

export const metadata: Metadata = {
  title: "Austin Pinball Collective: Sign In",
  description: "Sign in page for the Austin Pinball Collective.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="bg-zinc-700">
        <div id="page-root" className="m-1 md:m-2">
          {children}
        </div>
      </body>
    </html>
  );
}
