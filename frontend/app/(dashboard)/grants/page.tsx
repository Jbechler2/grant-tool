'use client'

import GrantCard from '@/app/components/features/grants/GrantCard';
import { useGrants } from '@/app/hooks/useGrants'
import { Grant } from '@/types';
import { ReactNode } from 'react';

export default function GrantsPage() {
  const { data: grants, isLoading, error } = useGrants()

  let grantChildren = [] as ReactNode
  
  if (grants && grants.length > 0){
    grantChildren = grants.map((grant: Grant)  => {
      return (
        GrantCard(grant)
      )
    });
  }

  if (isLoading) return <p>Loading...</p> 
  if (error) return <p>Something went wrong</p>
    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900">Grants</h1>
            <div>
              <pre>{grantChildren}</pre>
            </div>
        </div>
    )
}