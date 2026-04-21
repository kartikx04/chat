'use client'

import { useEffect, useRef, useState } from 'react'
import { cn, formatTime, getInitials } from '@/lib/utils'
import { ChatSkeleton } from '@/components/ui/Skeleton'
import type { Chat, ContactList, User } from '@/types'

interface ChatWindowProps {
  user: User
  contact: ContactList | null
  messages: Chat[]
  onSend: (message: string) => void
  connected: boolean
  historyLoading?: boolean
}

export function ChatWindow({ user, contact, messages, onSend, connected, historyLoading }: ChatWindowProps) {
  const [input, setInput] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    const trimmed = input.trim()
    if (!trimmed || !contact) return
    onSend(trimmed)
    setInput('')
    inputRef.current?.focus()
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  // Empty state — no contact selected
  if (!contact) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center gap-4 bg-background">
        <div className="w-16 h-16 rounded-2xl bg-card border border-border flex items-center justify-center">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className="text-muted-foreground">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
        </div>
        <div className="text-center">
          <p className="text-sm font-medium text-foreground">No conversation open</p>
          <p className="text-xs text-muted-foreground mt-1">Select a contact to start chatting</p>
        </div>
      </div>
    )
  }

  // Group messages by date
  const grouped: { date: string; messages: Chat[] }[] = []
  messages.forEach(msg => {
    const date = new Date(msg.created_at_unix * 1000).toDateString()
    const last = grouped[grouped.length - 1]
    if (last && last.date === date) {
      last.messages.push(msg)
    } else {
      grouped.push({ date, messages: [msg] })
    }
  })

  return (
    <div className="flex-1 flex flex-col bg-background min-w-0">
      {/* Chat header */}
      <div className="flex items-center gap-3 px-6 py-4 border-b border-border bg-card/50">
        <div className="w-9 h-9 rounded-full bg-primary/20 text-primary flex items-center justify-center text-xs font-bold font-display">
          {getInitials(contact.username)}
        </div>
        <div>
          <p className="text-sm font-medium text-foreground">{contact.username}</p>
          <div className="flex items-center gap-1.5">
            <span className={cn(
              'w-1.5 h-1.5 rounded-full',
              connected ? 'bg-primary animate-pulse-dot' : 'bg-muted-foreground'
            )} />
            <span className="text-xs text-muted-foreground">
              {connected ? 'connected' : 'reconnecting...'}
            </span>
          </div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-6 py-4 space-y-6">
        {historyLoading && <ChatSkeleton />}
        {messages.length === 0 && (
          <div className="flex flex-col items-center justify-center h-full gap-2 text-center">
            <p className="text-sm text-muted-foreground">No messages yet.</p>
            <p className="text-xs text-muted-foreground">Say hi to {contact.username}!</p>
          </div>
        )}

        {grouped.map(group => (
          <div key={group.date}>
            {/* Date divider */}
            <div className="flex items-center gap-3 mb-4">
              <div className="flex-1 h-px bg-border" />
              <span className="text-xs text-muted-foreground px-2">
                {new Date(group.date).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric' })}
              </span>
              <div className="flex-1 h-px bg-border" />
            </div>

            <div className="space-y-1.5">
              {group.messages.map((msg, i) => {
                const isSelf = msg.is_self || msg.from_id === user.id
                const prevMsg = group.messages[i - 1]
                const isFirstInGroup = !prevMsg || (prevMsg.from_id !== msg.from_id)

                return (
                  <div
                    key={msg.id}
                    className={cn(
                      'flex items-end gap-2',
                      isSelf ? 'justify-end' : 'justify-start',
                      isFirstInGroup ? 'mt-3' : 'mt-0.5'
                    )}
                  >
                    {/* Avatar for others */}
                    {!isSelf && isFirstInGroup && (
                      <div className="w-6 h-6 rounded-full bg-secondary text-xs font-bold font-display flex items-center justify-center text-muted-foreground shrink-0 mb-0.5">
                        {getInitials(contact.username)}
                      </div>
                    )}
                    {!isSelf && !isFirstInGroup && <div className="w-6 shrink-0" />}

                    <div className={cn('group flex flex-col max-w-[65%]', isSelf && 'items-end')}>
                      <div
                        className={cn(
                          'px-3.5 py-2 text-sm leading-relaxed',
                          isSelf
                            ? 'bg-primary text-primary-foreground rounded-2xl rounded-br-sm'
                            : 'bg-card border border-border text-foreground rounded-2xl rounded-bl-sm'
                        )}
                      >
                        {msg.message}
                      </div>
                      <span className="text-[10px] text-muted-foreground mt-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        {formatTime(msg.created_at_unix)}
                      </span>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div className="px-6 py-4 border-t border-border bg-card/30">
        <div className="flex items-end gap-3 bg-secondary rounded-xl px-4 py-3 border border-border focus-within:border-primary/30 transition-colors">
          <textarea
            ref={inputRef}
            rows={1}
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder={`Message ${contact.username}...`}
            className="flex-1 bg-transparent text-sm text-foreground placeholder:text-muted-foreground resize-none focus:outline-none leading-relaxed max-h-32"
            style={{ minHeight: '20px' }}
          />
          <button
            onClick={handleSend}
            disabled={!input.trim() || !connected}
            className={cn(
              'shrink-0 w-8 h-8 rounded-lg flex items-center justify-center transition-all',
              input.trim() && connected
                ? 'bg-primary text-primary-foreground hover:bg-primary/90 hover:scale-105'
                : 'bg-border text-muted-foreground cursor-not-allowed'
            )}
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
              <line x1="22" y1="2" x2="11" y2="13"/>
              <polygon points="22 2 15 22 11 13 2 9 22 2"/>
            </svg>
          </button>
        </div>
        <p className="text-[10px] text-muted-foreground mt-2 text-center">
          Enter to send · Shift+Enter for new line
        </p>
      </div>
    </div>
  )
}