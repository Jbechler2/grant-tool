'use client'

import { useClients } from "@/app/hooks/useClients"
import { Client } from "@/types"
import { ReactNode, useState } from "react"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { Label } from '@/components/ui/label'
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import apiClient from "@/lib/api"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import ClientsPage from "../page"

const clientSchema = z.object({
  name: z.string().min(1, "Funder name is required"),
  contact_name: z.string().min(1, "Contact name is required"),
  contact_phone: z.string().regex(/^\d+$/, "Must contain digits only").optional().or(z.literal("")),
  contact_email: z.string().email("Invalid email address").optional().or(z.literal("")),
  notes: z.string().optional(),

})

type ClientFormData = z.infer<typeof clientSchema>

export default function AddClients() {
  const router = useRouter()
  const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: {errors, isSubmitting },
  } = useForm<ClientFormData>({
    resolver: zodResolver(clientSchema)
  })

  const onSubmit = async(data: ClientFormData) => {
    setServerError(null)
    try{
      await apiClient.post('/clients', data)
      router.refresh()
      router.push('/clients')
    } catch {
      setServerError('Something went wrong.')
    }
  }

  return (
  <div>
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="space-y-2">
        <Label>Client Name</Label>
        <Input 
          id="clientName"
          type="text"
          {...register('name')}
        />
        {errors.contact_name && (
          <p className="text-sm text-red-500">{errors.contact_name.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label>Contact Name</Label>
        <Input 
          id="contactName"
          type="text"
          {...register('contact_name')}
        />
        {errors.contact_name && (
          <p className="text-sm text-red-500">{errors.contact_name.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label>Contact Phone</Label>
        <Input 
          id="contactPhone"
          type=""
          {...register('contact_phone')}
        />
        {errors.contact_phone && (
          <p className="text-sm text-red-500">{errors.contact_phone.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label>Contact Email</Label>
        <Input 
          id=""
          type="text"
          {...register('contact_email')}
        />
        {errors.contact_email && (
          <p className="text-sm text-red-500">{errors.contact_email.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label>Notes</Label>
        <Input 
          id="notes"
          type="text"
          {...register('notes')}
        />
        {errors.notes && (
          <p className="text-sm text-red-500">{errors.notes.message}</p>
        )}
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