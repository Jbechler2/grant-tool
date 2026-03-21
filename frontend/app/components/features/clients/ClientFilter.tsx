'use client'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useEffect, useState } from 'react'
import { useDebounce } from '../grants/GrantFilter'

interface ClientFilterBarProps {
  onSearchChange: (value: string) => void
}

export default function ClientFilter({ onSearchChange }: ClientFilterBarProps) {
  const [searchText, setSearchText] = useState("")
  const debouncedSearch = useDebounce(searchText, 300)

  useEffect(() => {
    onSearchChange(debouncedSearch)
  }, [debouncedSearch])

  return (
    <div className='bg-gray-100 p-5 rounded-lg'>
      <Label htmlFor='searchText'>Search</Label>
      <Input
        id="searchText"
        type="text"
        onChange={(e) => setSearchText(e.target.value)}
      />
    </div>
  )
}