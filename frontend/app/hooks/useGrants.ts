import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { Grant } from '@/types'

export function useGrants() {
  return useQuery({
    queryKey: ['grants'],
    queryFn: async () => {
      const response = await apiClient.get<Grant[]>('/grants')
      return response.data
    }
  })
}