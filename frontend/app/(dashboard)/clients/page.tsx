'use client'

import ClientCard from "@/app/components/features/clients/ClientCard";
import { useClients } from "@/app/hooks/useClients"
import { Client } from "@/types";
import { ReactNode } from "react";

export default function ClientsPage() {
  const { data: clients, isLoading, error } = useClients()

  let clientChildren = [] as ReactNode
  
      if (clients && clients.length > 0){
        clientChildren = clients.map((app: Client)  => {
          return (
            ClientCard(app)
          )
        });
      }

  if (isLoading) return <p>Loading...</p> 
  if (error) return <p>Something went wrong</p>
    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900">Clients</h1>
            <div>
              <pre>{clientChildren}</pre>
            </div>
        </div>
    )
}