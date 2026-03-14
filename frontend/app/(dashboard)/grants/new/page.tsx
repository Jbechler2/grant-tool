'use client'

import { useGrants } from "@/app/hooks/useGrants"
import { Grant } from "@/types"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { ReactNode, useState } from "react"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { Label } from '@/components/ui/label'
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import apiClient from "@/lib/api"

const grantSchema = z.object({
  title: z.string('Please enter a grant name').min(1, "Grant Title is required"),
  funder_name: z.string('Please enter a funder').min(1, "Funder is required"), 
  funder_website: z.string().optional(),
  description: z.string().optional(),
  award_amount_min: z.number().positive().optional().or(z.nan().transform(() => undefined)),
  award_amount_max: z.number().positive().optional().or(z.nan().transform(() => undefined)),
  eligibility_notes: z.string().optional(),
  estimated_application_hours: z.number().positive().optional().or(z.nan().transform(() => undefined)),
})

type GrantFormData = z.infer<typeof grantSchema>

export default function GrantsPage() {
const router = useRouter()
const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: {errors, isSubmitting },
  } = useForm<GrantFormData>({
    resolver: zodResolver(grantSchema)
  })

const onSubmit = async (data: GrantFormData) => {
  setServerError(null)
  try{
    await apiClient.post('/grants', data)
        router.push('/grants')
        router.refresh()
  } catch {
    setServerError('Something went wrong.')
  }
}

  return (
  <div>
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="space-y-2">
        <Label htmlFor="grantName">Grant Name*</Label>
        <Input 
          id="grantName"
          type="text"
          placeholder="Example Grant Name"
          {...register('title')}
        />
        {errors.title && (
          <p className="text-sm text-red-500">{errors.title.message}</p>
        )}
      </div>
      <div className="space-y-2">
      <Label htmlFor="funderName">Funder Name*</Label>
        <Input 
          id="funderName"
          type="text"
          placeholder="Example Funder Name"
          {...register('funder_name')}
        />
        {errors.funder_name && (
          <p className="text-sm text-red-500">{errors.funder_name.message}</p>
        )}
      </div>
      <Label htmlFor="funderWebsite">Funder Website</Label>
      <Input 
        id="funderWebsite"
        type="text"
        placeholder="Example Funder Website"
        {...register('funder_website')}
      />
      <Label htmlFor="description">Description</Label>
      <Input 
        id="description"
        type="text"
        placeholder="Grant Description..."
        {...register('description')}
      />
      <div className="space-y-2">
        <Label htmlFor="minAward">Minimum Award</Label>
        <Input 
          id="minAward"
          type="number"
          {...register('award_amount_min', { valueAsNumber: true})}
        />
        {errors.award_amount_min && (
            <p className="text-sm text-red-500">{errors.award_amount_min.message}</p>
        )}
      </div>
      <div className="space-y-2">
      <Label htmlFor="maxAward">Maximum Award</Label>
        <Input 
          id="maxAward"
          type="number"
          {...register('award_amount_max', { valueAsNumber: true})}
        />
        {errors.award_amount_max && (
          <p className="text-sm text-red-500">{errors.award_amount_max.message}</p>
        )}
      </div>
      <div className="space-y-2">
      <Label htmlFor="estimatedApplicationHours">Est. Application Hours</Label>
      <Input 
        id="estimatedApplicationHours"
        type="number"
        {...register('estimated_application_hours', { valueAsNumber: true})}
      />
      </div>
      <Button
        type="submit"
        className="w-full"
        disabled={isSubmitting}
      >
        {isSubmitting ? 'Submitting...' : 'Submit'}
      </Button>
    </form>
  </div>
  )
}