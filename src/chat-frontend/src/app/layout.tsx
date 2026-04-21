import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Kayee — Real-time Chat',
  description: 'Fast, simple real-time messaging',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="dark">
      <body className="noise-bg min-h-screen">
        {children}
      </body>
    </html>
  )
}