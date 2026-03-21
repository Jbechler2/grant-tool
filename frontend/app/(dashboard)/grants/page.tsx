'use client'

import { useGrants } from '@/app/hooks/useGrants'
import { Grant } from '@/types'
import { useRouter } from 'next/navigation'
import GrantFilter from '@/app/components/features/grants/GrantFilter'
import { useEffect, useMemo, useState } from 'react'
import GrantListView from '@/app/components/features/grants/GrantListView'
import Fuse from 'fuse.js'

export default function GrantsPage() {
  const { data: grants = [], isLoading } = useGrants()
  const [searchQuery, setSearchQuery] = useState("")

  const fuse = useMemo(() => new Fuse(grants, {
    keys: ["title", "funder_name"],
    threshold: 0.3,
  }), [grants])

  const filteredGrants = useMemo(() => {
    if (!searchQuery) return grants
    return fuse.search(searchQuery).map(result => result.item)
  }, [searchQuery, fuse])

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="bg-muted rounded-lg h-32 animate-pulse" />
        ))}
      </div>
    )
  }

  if (grants.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <p className="text-sm font-medium text-foreground">No grants yet</p>
        <p className="text-xs text-muted-foreground mt-1">Add your first grant to get started</p>
      </div>
    )
  }

  return (
    <div className='flex flex-col'>
      <div className='w-1/2 mb-10'>
        <GrantFilter onSearchChange={setSearchQuery}></GrantFilter>
      </div>
      <div>
        <GrantListView ViewMode='card' Grants={filteredGrants}></GrantListView>
      </div>
    </div>
    
  )
}