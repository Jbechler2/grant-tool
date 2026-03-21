'use client'

import ClientFilter from '@/app/components/features/clients/ClientFilter'
import { useClients } from '@/app/hooks/useClients'
import { Client } from '@/types'
import { useRouter } from 'next/navigation'
import { useMemo, useState } from 'react'
import Fuse from 'fuse.js'
import ClientListView from '@/app/components/features/clients/ClientListView'


export default function ClientsPage() {
  const { data: clients = [], isLoading } = useClients()
  const [searchQuery, setSearchQuery] = useState("")

  const fuse = useMemo(() => new Fuse(clients, {
    keys: ["name"],
    threshold: 0.3,
  }), [clients])

  const filteredClients = useMemo(() => {
    if (!searchQuery) return clients
    return fuse.search(searchQuery).map(result => result.item)
  }, [searchQuery, fuse])

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
    <div className='flex flex-col'>
      <div className='w-1/2 mb-10'>
        <ClientFilter onSearchChange={setSearchQuery}></ClientFilter>
      </div>
      <div>
        <ClientListView ViewMode='card' Clients={filteredClients}></ClientListView>
      </div>
    </div>
  )
}