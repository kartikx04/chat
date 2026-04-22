'use client'

import { useState } from 'react'
import { cn, formatTime, getInitials } from '@/lib/utils'
import { apiClient } from '@/lib/api'
import type { ContactList, User } from '@/types'

interface ContactSidebarProps {
  user: User
  contacts: ContactList[]
  activeContact: ContactList | null
  onSelectContact: (contact: ContactList) => void
  onContactsUpdate: (contacts: ContactList[]) => void
  collapsed: boolean
  onToggleCollapse: () => void
}

export function ContactSidebar({
  user,
  contacts,
  activeContact,
  onSelectContact,
  onContactsUpdate,
  collapsed,
  onToggleCollapse,
}: ContactSidebarProps) {
  const [addOpen, setAddOpen] = useState(false)
  const [searchVal, setSearchVal] = useState('')
  const [addError, setAddError] = useState('')
  const [adding, setAdding] = useState(false)

  const handleAddContact = async () => {
    if (!searchVal.trim()) return
    setAdding(true)
    setAddError('')
    try {
      const verify = await apiClient.verifyContact(searchVal.trim())
      if (!verify.data.status) {
        setAddError('User not found')
        return
      }
      
      const res = await apiClient.addContact(user.id, searchVal.trim())
      console.log('addContact response:', res.data)
      if (!res.data.status) {
        setAddError(res.data.message || 'Failed to add contact')
        return
      }
      
      const updated = await apiClient.getContacts(user.id)
      console.log('getContacts response:', updated.data)
      if (updated.data.data) {
        console.log('Updating contacts with:', updated.data.data)
        onContactsUpdate(updated.data.data)
      }
      
      setSearchVal('')
      setAddOpen(false)
    } catch (err) {
      console.error('addContact error:', err)
      setAddError('Something went wrong')
    } finally {
      setAdding(false)
    }
  }

  return (
    <aside
      className={cn(
        'flex flex-col border-r border-border bg-card transition-all duration-300 ease-in-out',
        collapsed ? 'w-16' : 'w-72'
      )}
    >
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-4 border-b border-border min-h-[60px]">
        {!collapsed && (
          <span className="font-display font-semibold text-sm text-foreground">Messages</span>
        )}
        <button
          onClick={onToggleCollapse}
          className="ml-auto p-1.5 rounded-lg hover:bg-secondary text-muted-foreground hover:text-foreground transition-colors"
          title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            {collapsed
              ? <><path d="M9 18l6-6-6-6"/></>
              : <><path d="M15 18l-6-6 6-6"/></>
            }
          </svg>
        </button>
      </div>

      {/* Add contact */}
      {!collapsed && (
        <div className="px-3 py-3 border-b border-border">
          {addOpen ? (
            <div className="space-y-2">
              <input
                type="text"
                placeholder="Username to add..."
                value={searchVal}
                onChange={e => { setSearchVal(e.target.value); setAddError('') }}
                onKeyDown={e => e.key === 'Enter' && handleAddContact()}
                className="w-full px-3 py-2 text-sm bg-secondary rounded-lg border border-border focus:outline-none focus:border-primary/50 placeholder:text-muted-foreground"
                autoFocus
              />
              {addError && <p className="text-xs text-red-400">{addError}</p>}
              <div className="flex gap-2">
                <button
                  onClick={handleAddContact}
                  disabled={adding}
                  className="flex-1 py-1.5 text-xs bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50"
                >
                  {adding ? '...' : 'Add'}
                </button>
                <button
                  onClick={() => { setAddOpen(false); setSearchVal(''); setAddError('') }}
                  className="flex-1 py-1.5 text-xs bg-secondary text-secondary-foreground rounded-lg hover:bg-border transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          ) : (
            <button
              onClick={() => setAddOpen(true)}
              className="w-full flex items-center gap-2 px-3 py-2 text-xs text-muted-foreground hover:text-foreground hover:bg-secondary rounded-lg transition-colors"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              Add contact
            </button>
          )}
        </div>
      )}

      {/* Contact list */}
      <div className="flex-1 overflow-y-auto py-2">
        {contacts.length === 0 && !collapsed && (
          <div className="px-4 py-8 text-center">
            <p className="text-xs text-muted-foreground">No contacts yet.</p>
            <p className="text-xs text-muted-foreground mt-1">Add someone to start chatting.</p>
          </div>
        )}
        {contacts.map(contact => (
          <button
            key={contact.id}
            onClick={() => onSelectContact(contact)}
            className={cn(
              'w-full flex items-center gap-3 px-3 py-2.5 transition-colors hover:bg-secondary/50 group',
              activeContact?.id === contact.id && 'bg-secondary'
            )}
          >
            {/* Avatar */}
            <div className={cn(
              'shrink-0 w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold font-display',
              activeContact?.id === contact.id
                ? 'bg-primary text-primary-foreground'
                : 'bg-secondary text-muted-foreground group-hover:bg-border'
            )}>
              {getInitials(contact.username)}
            </div>
            {/* Info */}
            {!collapsed && (
              <div className="flex-1 min-w-0 text-left">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-foreground truncate">{contact.username}</span>
                  {contact.last_activity > 0 && (
                    <span className="text-xs text-muted-foreground shrink-0 ml-2">
                      {formatTime(contact.last_activity)}
                    </span>
                  )}
                </div>
              </div>
            )}
          </button>
        ))}
      </div>

      {/* User info at bottom */}
      <div className={cn(
        'border-t border-border px-3 py-3 flex items-center gap-3',
        collapsed && 'justify-center'
      )}>
        <div className="w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-xs font-bold font-display shrink-0">
          {getInitials(user.username)}
        </div>
        {!collapsed && (
          <div className="flex-1 min-w-0">
            <p className="text-xs font-medium text-foreground truncate">{user.username}</p>
            <p className="text-xs text-muted-foreground truncate">{user.email}</p>
          </div>
        )}
      </div>
    </aside>
  )
}