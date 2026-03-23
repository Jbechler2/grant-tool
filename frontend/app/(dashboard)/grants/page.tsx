'use client'

import { useGrants } from '@/app/hooks/useGrants'
import { Grant } from '@/types'
import { useRouter } from 'next/navigation'
import GrantFilter from '@/app/components/features/grants/GrantFilter'
import { useEffect, useMemo, useState } from 'react'
import GrantListView from '@/app/components/features/grants/GrantListView'
import Fuse from 'fuse.js'
import { Button } from '@/components/ui/button'
import { List, Square } from 'lucide-react'

export default function GrantsPage() {
  const { data: grants = [], isLoading } = useGrants()
  const [searchQuery, setSearchQuery] = useState("")
  const [viewMode, setViewMode] = useState("card")

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
          <div className='lg:w-1/2 sm:w-4/5 flex mb-1 flex-row justify-between'>
            <div className='lg:w-3/5 sm:w-1/5'>
              <GrantFilter onSearchChange={setSearchQuery}></GrantFilter>
            </div>
            <div className='flex flex-col lg:w-1/5 sm:w-full justify-between rounded-lg self-end'>
              <div className='flex justify-between'>
                <Button
                  variant='ghost'
                  className='bg-white-400'
                  onClick={() => setViewMode('list')}
                >
                  <List />
                </Button>
                <div className='my-auto mx-2'>|</div>
                <Button
                  variant='ghost'
                  className='bg-white-400'
                  onClick={() => setViewMode('card')}
                ><Square />
                </Button>
              </div>
            </div>
          </div>
          <div>
            <GrantListView ViewMode={viewMode} Grants={filteredGrants}></GrantListView>
          </div>
        </div>
  )
}