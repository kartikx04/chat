import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatTime(unixTs: number): string {
  const date = new Date(unixTs * 1000)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
 
export function formatDate(unixTs: number): string {
  const date = new Date(unixTs * 1000)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
 
  if (days === 0) return 'Today'
  if (days === 1) return 'Yesterday'
  return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
}
 
export function getInitials(username: string): string {
  return username.slice(0, 2).toUpperCase()
}
 