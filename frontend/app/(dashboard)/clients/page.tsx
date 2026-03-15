'use client'

import { useClients } from '@/app/hooks/useClients'
import { Client } from '@/types'
import { useRouter } from 'next/navigation'

function getInitials(name: string) {
  return name
    .split(' ')
    .map(word => word[0])
    .join('')
    .slice(0, 2)
    .toUpperCase()
}

const avatarColors = [
  'bg-blue-100 text-blue-800',
  'bg-teal-100 text-teal-800',
  'bg-orange-100 text-orange-800',
  'bg-purple-100 text-purple-800',
  'bg-amber-100 text-amber-800',
]

function getAvatarColor(name: string) {
  const index = name.charCodeAt(0) % avatarColors.length
  return avatarColors[index]
}

function ClientCard({ client }: { client: Client }) {
  const router = useRouter()

  return (
    <div
      onClick={() => router.push(`/clients/${client.id}`)}
      className="bg-white border border-border rounded-lg p-6 flex flex-col items-center text-center gap-3 cursor-pointer hover:border-border/80 transition-colors"
    >
      <div className={`w-12 h-12 rounded-full flex items-center justify-center text-sm font-medium ${getAvatarColor(client.name)}`}>
        {getInitials(client.name)}
      </div>
      <p className="text-sm font-medium text-foreground leading-tight">{client.name}</p>
    </div>
  )
}

export default function ClientsPage() {
  const { data: clients = [], isLoading } = useClients()

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
        {Array.from({ length: 5 }).map((_, i) => (
          <div key={i} className="bg-muted rounded-lg h-32 animate-pulse" />
        ))}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
      {clients.map(client => (
        <ClientCard key={client.id} client={client} />
      ))}
    </div>
  )
}