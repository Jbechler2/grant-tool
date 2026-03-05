import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { Client } from '@/types'

export function useClients() {
  return useQuery({
    queryKey: ['clients'],
    queryFn: async () => {
      const response = await apiClient.get<Client[]>('/clients')
      return response.data
    },
  })
}