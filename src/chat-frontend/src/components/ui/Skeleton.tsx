'use client'

import { cn } from '@/lib/utils'

function Skeleton({ className, style }: { className?: string; style?: React.CSSProperties }) {
  return (
    <div className={cn('animate-pulse rounded-lg bg-secondary', className)} style={style} />
  )
}

export function ChatSkeleton() {
  return (
    <div className="flex-1 flex flex-col px-6 py-4 space-y-4 overflow-hidden">
      {/* Incoming */}
      <div className="flex items-end gap-2">
        <Skeleton className="w-6 h-6 rounded-full shrink-0" />
        <Skeleton className="h-9 w-48 rounded-2xl rounded-bl-sm" />
      </div>
      {/* Outgoing */}
      <div className="flex justify-end">
        <Skeleton className="h-9 w-36 rounded-2xl rounded-br-sm" />
      </div>
      {/* Incoming long */}
      <div className="flex items-end gap-2">
        <Skeleton className="w-6 h-6 rounded-full shrink-0" />
        <Skeleton className="h-9 w-64 rounded-2xl rounded-bl-sm" />
      </div>
      {/* Outgoing short */}
      <div className="flex justify-end">
        <Skeleton className="h-9 w-20 rounded-2xl rounded-br-sm" />
      </div>
      {/* Incoming */}
      <div className="flex items-end gap-2">
        <Skeleton className="w-6 h-6 rounded-full shrink-0" />
        <Skeleton className="h-9 w-40 rounded-2xl rounded-bl-sm" />
      </div>
    </div>
  )
}

export function ContactSkeleton() {
  return (
    <div className="flex flex-col gap-1 px-3 py-2">
      {[64, 48, 56, 40, 52].map((w, i) => (
        <div key={i} className="flex items-center gap-3 px-2 py-2">
          <Skeleton className="w-9 h-9 rounded-full shrink-0" />
          <div className="flex-1 space-y-1.5">
            <Skeleton className={`h-3 rounded`} style={{ width: `${w}%` }} />
            <Skeleton className="h-2.5 rounded w-1/2" />
          </div>
        </div>
      ))}
    </div>
  )
}