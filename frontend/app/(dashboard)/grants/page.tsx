'use client'

import { useGrants } from '@/app/hooks/useGrants'
import { Grant } from '@/types'
import { useRouter } from 'next/navigation'

function GrantCard({ grant }: { grant: Grant }) {
  const router = useRouter()

  return (
    <div
      onClick={() => router.push(`/grants/${grant.id}`)}
      className="bg-white border border-border rounded-lg p-5 flex flex-col gap-2 cursor-pointer hover:border-border/80 transition-colors"
    >
      <div className="flex justify-between items-center">
        <p className="text-xs text-muted-foreground">{grant.funder_name}</p>
        {grant.award_amount_max && (
          <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-green-100 text-green-800">
            ${Number(grant.award_amount_max).toLocaleString()}
          </span>
        )}
      </div>
      <p className="text-sm font-medium text-foreground leading-snug">{grant.title}</p>
      {grant.description && (
        <p className="text-xs text-muted-foreground leading-relaxed line-clamp-2">
          {grant.description}
        </p>
      )}
    </div>
  )
}

export default function GrantsPage() {
  const { data: grants = [], isLoading } = useGrants()

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
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {grants.map(grant => (
        <GrantCard key={grant.id} grant={grant} />
      ))}
    </div>
  )
}