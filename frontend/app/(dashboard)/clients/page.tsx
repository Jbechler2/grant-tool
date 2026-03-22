'use client'

import ClientFilter from '@/app/components/features/clients/ClientFilter'
import { useClients } from '@/app/hooks/useClients'
import { Client } from '@/types'
import { useRouter } from 'next/navigation'
import { useMemo, useState } from 'react'
import Fuse from 'fuse.js'
import ClientListView from '@/app/components/features/clients/ClientListView'
import { Button } from '@/components/ui/button'
import { Square, List } from 'lucide-react'

export default function ClientsPage() {
  const { data: clients = [], isLoading } = useClients()
  const [searchQuery, setSearchQuery] = useState("")
  const [viewMode, setViewMode] = useState("card")

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
      <div className='lg:w-1/2 sm:w-4/5 mb-10 flex flex-row justify-between'>
        <div className='lg:w-3/5 sm:w-1/5'>
          <ClientFilter onSearchChange={setSearchQuery}></ClientFilter>
        </div>
        <div className='flex flex-col m-3 lg:w-1/5 sm:w-full bg-blue-200 justify-between rounded-lg p-3'>
          <div className='flex justify-center'>
            <h1>View Mode</h1>
          </div>
          <div className='flex justify-center'>
            <Button
              className='bg-gray-400 mr-5'
              onClick={() => setViewMode('list')}
            >
              <List />
            </Button>
            <Button
              className='bg-gray-400'
              onClick={() => setViewMode('card')}
            ><Square />
            </Button>
          </div>
        </div>
      </div>
      <div>
        <ClientListView ViewMode={viewMode} Clients={filteredClients}></ClientListView>
      </div>
    </div>
  )
}