'use client'

import { useClients } from "@/app/hooks/useClients"

export default function ClientsPage() {
  const { data: clients, isLoading, error } = useClients()

  if (isLoading) return <p>Loading...</p> 
  if (error) return <p>Something went wrong</p>
    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900">Clients</h1>
            <pre>{JSON.stringify(clients, null, 2)}</pre>
        </div>
    )
}