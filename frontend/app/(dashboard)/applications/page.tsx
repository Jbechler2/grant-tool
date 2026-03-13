'use client'

import ApplicationCard from "@/app/components/features/applications/ApplicationCard"
import { useApplications } from "@/app/hooks/useApplications"
import { Card, CardTitle } from "@/components/ui/card"
import { Application } from "@/types"
import Link from "next/link"
import { ReactNode } from "react"

export default function ApplicationsPage() {
   const { data: applications, isLoading, error } = useApplications()
  
    if (isLoading) return <p>Loading...</p> 
    if (error) return <p>Something went wrong</p>
    let appChildren = [] as ReactNode

    if (applications && applications.length > 0){
      appChildren = applications.map((app: Application)  => {
        return (
          ApplicationCard(app)
        )
      });
    }
    
    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900">Applications</h1>
            <pre className="flex flex-row">{appChildren}</pre>
        </div>
    )
}