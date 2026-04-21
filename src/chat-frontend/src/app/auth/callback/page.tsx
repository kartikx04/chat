'use client'

import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useUser } from '@/hooks/useUser'

// Your Go /callback endpoint should redirect to this page with user info
// e.g. /auth/callback?id=<uuid>&username=<name>&email=<email>
// Adjust the param names to match what your Go callback actually sends

export default function CallbackPage() {
  const router = useRouter()
  const params = useSearchParams()
  const { setUser } = useUser()

  useEffect(() => {
    const id = params.get('id')
    const username = params.get('username')
    const email = params.get('email') || ''

    if (id && username) {
      setUser({ id, username, email })
      router.replace('/home')
    } else {
      // No user info — go back to landing
      router.replace('/')
    }
  }, [params, router, setUser])

  return (
    <div className="h-screen flex items-center justify-center">
      <div className="flex flex-col items-center gap-3">
        <div className="w-8 h-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
        <p className="text-sm text-muted-foreground">Signing you in...</p>
      </div>
    </div>
  )
}