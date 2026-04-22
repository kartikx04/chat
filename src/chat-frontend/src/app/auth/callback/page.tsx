'use client'

export const dynamic = 'force-dynamic'

import { Suspense } from 'react'
import { useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useUser } from '@/hooks/useUser'

function CallbackHandler() {
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

export default function CallbackPage() {
  return (
    <Suspense fallback={
      <div className="h-screen flex items-center justify-center">
        <div className="w-8 h-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
      </div>
    }>
      <CallbackHandler />
    </Suspense>
  )
}