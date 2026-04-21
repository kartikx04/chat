export interface User {
  id: string
  username: string
  email?: string
  avatar?: string
}

export interface Chat {
  id: string
  from_id: string
  to_id: string
  message: string
  created_at: string
  created_at_unix: number
  is_self: boolean
}

export interface ContactList {
  id: string
  username: string
  last_activity: number
}

export interface WSMessage {
  type: 'bootup' | 'message'
  user?: string
  user_id?: string
  chat?: {
    from_id: string
    to_id: string
    message: string
  }
}

export interface APIResponse<T = unknown> {
  status: boolean
  message?: string
  data?: T
  total?: number
}