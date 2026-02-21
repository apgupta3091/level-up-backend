import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import { ClerkProvider } from '@clerk/nextjs'
import './globals.css'

const geistSans = Geist({ variable: '--font-geist-sans', subsets: ['latin'] })
const geistMono = Geist_Mono({ variable: '--font-geist-mono', subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Level Up Backend',
  description: 'From mid-level to senior backend engineer in 6 months.',
}

const clerkKey = process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY ?? ''
const clerkConfigured = clerkKey.startsWith('pk_') && clerkKey !== 'pk_test_...'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const inner = (
    <html lang="en">
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        {children}
      </body>
    </html>
  )

  if (!clerkConfigured) return inner

  return <ClerkProvider>{inner}</ClerkProvider>
}
