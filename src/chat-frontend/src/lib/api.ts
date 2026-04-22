import axios from 'axios'
import type { APIResponse, Chat, ContactList } from '@/types'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  withCredentials: false, // sends session cookie automatically
})

export const apiClient = {
  // Health check
  health: () => api.get('/health'),

  // Verify a contact exists before adding
  verifyContact: (username: string) =>
  api.get<APIResponse>(`/verify-contact?username=${encodeURIComponent(username)}`),

  // Add contact (bidirectional on backend)
  addContact: (id: string, contact: string) =>
    api.get<APIResponse>(`/add-contact?id=${id}&contact=${contact}`),

  // Fetch contact list for a user
  getContacts: (id: string) =>
    api.get<APIResponse<ContactList[]>>(`/contacts?id=${id}`),

  // Fetch chat history between two users
  getChatHistory: (u1: string, u2: string, fromTs = '0', toTs = '+inf') =>
    api.get<APIResponse<Chat[]>>(
      `/chat-history?u1=${u1}&u2=${u2}&from-ts=${fromTs}&to-ts=${toTs}`
    ),
}

export default api