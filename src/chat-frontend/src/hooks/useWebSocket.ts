'use client'

import { useEffect, useRef, useCallback, useState } from 'react'
import type { Chat, User } from '@/types'

interface UseWebSocketOptions {
  user: User | null
  onMessage: (chat: Chat) => void
}

export function useWebSocket({ user, onMessage }: UseWebSocketOptions) {
  const wsRef = useRef<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)
  const reconnectTimeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const onMessageRef = useRef(onMessage)
  onMessageRef.current = onMessage

  const connect = useCallback(() => {
    if (!user) return
    if (wsRef.current?.readyState === WebSocket.OPEN) return

    const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8081'
    const ws = new WebSocket(`${wsUrl}/ws`)
    wsRef.current = ws

    ws.onopen = () => {
      setConnected(true)
      // Send bootup message to register this connection with user ID
      ws.send(JSON.stringify({
        type: 'bootup',
        user: user.username,
        user_id: user.id,
      }))
    }

    ws.onmessage = (event) => {
      try {
        const data: Chat = JSON.parse(event.data)
        onMessageRef.current(data)
      } catch (e) {
        console.error('WS parse error:', e)
      }
    }

    ws.onclose = () => {
      setConnected(false)
      // Reconnect after 3s
      reconnectTimeout.current = setTimeout(() => connect(), 3000)
    }

    ws.onerror = (err) => {
      console.error('WS error:', err)
      ws.close()
    }
  }, [user])

  useEffect(() => {
    connect()
    return () => {
      clearTimeout(reconnectTimeout.current)
      wsRef.current?.close()
    }
  }, [connect])

  const sendMessage = useCallback((fromId: string, toId: string, message: string) => {
    if (wsRef.current?.readyState !== WebSocket.OPEN) return false
    wsRef.current.send(JSON.stringify({
      type: 'message',
      chat: { from_id: fromId, to_id: toId, message },
    }))
    return true
  }, [])

  return { connected, sendMessage }
}