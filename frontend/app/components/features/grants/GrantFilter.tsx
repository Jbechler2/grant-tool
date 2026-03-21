'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useEffect, useState } from 'react'

interface GrantFilterBarProps {
  onSearchChange: (value: string) => void
}

export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value)

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedValue(value), delay)
    return () => clearTimeout(timer)
  }, [value, delay])

  return debouncedValue
}

export default function GrantFilter({ onSearchChange }: GrantFilterBarProps) {
  const [searchText, setSearchText] = useState("")
  const debouncedSearch = useDebounce(searchText, 300)

  useEffect(() => {
    onSearchChange(debouncedSearch)
  }, [debouncedSearch])
  

  return (
    <div className='bg-gray-100 p-5 rounded-lg'>
      <Label htmlFor="searchText">Search</Label>
      <Input
        id="searchText"
        type="text"
        onChange={(e) => setSearchText(e.target.value)}
      />
    </div>
  )
}