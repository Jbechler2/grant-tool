import { useQuery } from "@tanstack/react-query";
import apiClient from "@/lib/api";
import { Application } from '@/types'

export function useApplications() {
  return useQuery({
    queryKey: ['applications'],
    queryFn: async () => {
      const response = await apiClient.get<Application[]>('/applications')
      return response.data
    },
  })
}