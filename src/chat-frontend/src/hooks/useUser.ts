'use client'

import { useState, useEffect, useCallback } from 'react'
import type { User } from '@/types'

export function useUser() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const stored = localStorage.getItem('chat_user')
    if (stored) {
      try {
        setUser(JSON.parse(stored))
      } catch {
        localStorage.removeItem('chat_user')
      }
    }
    setLoading(false)
  }, [])

  const setAndStoreUser = useCallback((u: User) => {
    localStorage.setItem('chat_user', JSON.stringify(u))
    setUser(u)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('chat_user')
    setUser(null)
  }, [])

  return { user, loading, setUser: setAndStoreUser, logout }
}