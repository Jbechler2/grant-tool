import { cookies } from 'next/headers'

export async function getServerApiClient() {
  const cookieStore = await cookies()
  const token = cookieStore.get('token')?.value

  return {
    get: async (path: string) => {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}${path}`,
        {
          headers: {
            'Authorization': token ? `Bearer ${token}` : '',
            'Content-Type': 'application/json',
          },
          cache: 'no-store',
        }
      )
      if(!response.ok) throw new Error(`API error: ${response.status}`)
      return response.json()
    }
  }
}