'use client'

import { useEffect, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import { useUser } from '@/hooks/useUser'
import { useWebSocket } from '@/hooks/useWebSocket'
import { apiClient } from '@/lib/api'
import { ContactSidebar } from '@/components/chat/ContactSidebar'
import { ChatWindow } from '@/components/chat/ChatWindow'
import type { Chat, ContactList } from '@/types'

export default function HomePage() {
  const router = useRouter()
  const { user, loading, logout } = useUser()
  const [contacts, setContacts] = useState<ContactList[]>([])
  const [activeContact, setActiveContact] = useState<ContactList | null>(null)
  const [messages, setMessages] = useState<Map<string, Chat[]>>(new Map())
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)
  const [contactsLoading, setContactsLoading] = useState(false)
  const [historyLoading, setHistoryLoading] = useState(false)

  // Redirect to landing if not logged in
  useEffect(() => {
    if (!loading && !user) {
      router.replace('/')
    }
  }, [user, loading, router])

  // Load contacts on mount
  useEffect(() => {
    if (!user) return
    setContactsLoading(true)
    apiClient.getContacts(user.id)
      .then(res => {
        if (res.data.data) setContacts(res.data.data)
      })
      .catch(console.error)
      .finally(() => setContactsLoading(false))
  }, [user])

  // Load chat history when contact is selected
  useEffect(() => {
    if (!user || !activeContact) return

    // If we already loaded history for this contact, skip
    if (messages.has(activeContact.id)) return

    setHistoryLoading(true)
    apiClient.getChatHistory(user.id, activeContact.id)
      .then(res => {
        const chats = res.data.data || []
        setMessages(prev => new Map(prev).set(activeContact.id, [...chats].reverse()))
      })
      .catch(console.error)
      .finally(() => setHistoryLoading(false))
  }, [activeContact, user, messages])

  // Handle incoming WebSocket messages
  const handleIncomingMessage = useCallback((chat: Chat) => {
    const contactId = chat.is_self ? chat.to_id : chat.from_id
    setMessages(prev => {
      const existing = prev.get(contactId) || []
      return new Map(prev).set(contactId, [...existing, chat])
    })

    // Update contact list order (bump to top)
    setContacts(prev => {
      const idx = prev.findIndex(c => c.id === contactId)
      if (idx === -1) return prev
      const updated = [...prev]
      const [contact] = updated.splice(idx, 1)
      return [{ ...contact, last_activity: chat.created_at_unix }, ...updated]
    })
  }, [])

  const { connected, sendMessage } = useWebSocket({
    user,
    onMessage: handleIncomingMessage,
  })

  const handleSend = (message: string) => {
    if (!user || !activeContact) return
    const sent = sendMessage(user.id, activeContact.id, message)
    if (!sent) console.warn('WebSocket not connected, message not sent')
  }

  const handleSelectContact = (contact: ContactList) => {
    setActiveContact(contact)
    // Auto-expand sidebar on mobile when selecting a contact
    if (window.innerWidth < 640) setSidebarCollapsed(true)
  }

  const currentMessages = activeContact ? (messages.get(activeContact.id) || []) : []

  if (loading || !user) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 rounded-full border-2 border-primary border-t-transparent animate-spin" />
          <p className="text-sm text-muted-foreground">Loading...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Top bar */}
      <header className="h-[60px] border-b border-border bg-card/50 flex items-center justify-between px-6 shrink-0">
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 rounded-md bg-primary flex items-center justify-center">
            <span className="text-primary-foreground font-display font-bold text-xs">K</span>
          </div>
          <span className="font-display font-semibold text-sm">Kayee</span>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5">
            <span className={`w-1.5 h-1.5 rounded-full ${connected ? 'bg-primary animate-pulse-dot' : 'bg-red-400'}`} />
            <span className="text-xs text-muted-foreground hidden sm:inline">
              {connected ? 'live' : 'offline'}
            </span>
          </div>
          <button
            onClick={logout}
            className="text-xs text-muted-foreground hover:text-foreground transition-colors px-2 py-1 rounded hover:bg-secondary"
          >
            Sign out
          </button>
        </div>
      </header>

      {/* Main layout */}
      <div className="flex flex-1 overflow-hidden">
        <ContactSidebar
          user={user}
          contacts={contacts}
          activeContact={activeContact}
          onSelectContact={handleSelectContact}
          onContactsUpdate={setContacts}
          collapsed={sidebarCollapsed}
          onToggleCollapse={() => setSidebarCollapsed(v => !v)}
        />
        <ChatWindow
          user={user}
          contact={activeContact}
          messages={currentMessages}
          onSend={handleSend}
          connected={connected}
          historyLoading={historyLoading}
        />
      </div>
    </div>
  )
}