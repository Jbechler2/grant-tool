'use client'

import { useApplications } from "@/app/hooks/useApplications"
import { useGrants } from "@/app/hooks/useGrants"
import { useClients } from "@/app/hooks/useClients"
import { Application, Client, Grant } from "@/types"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { ReactNode, useState } from "react"
import { Controller, useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { Label } from '@/components/ui/label'
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Combobox, ComboboxContent, ComboboxEmpty, ComboboxInput, ComboboxItem, ComboboxList } from "@/components/ui/combobox"
import { Select, SelectContent, SelectItem, SelectTrigger,SelectValue } from "@/components/ui/select"
import apiClient from "@/lib/api"

const applicationSchema = z.object({
  grant_id: z.string(),
  client_id: z.string(),
  title: z.string('Please enter an application name').min(1, 'Application name is required'),
  status: z.string(),
  notes: z.string().optional()
})

type ApplicationFormData = z.infer<typeof applicationSchema>

export default function AddApplication() {
const router = useRouter()
const [serverError, setServerError] = useState<string | null>(null)

const { data: grants, isLoading: grantsLoading } = useGrants()
const { data: clients, isLoading: clientsLoading } = useClients()

const {
  control,
  register,
  setValue,
  handleSubmit,
  formState: {errors, isSubmitting },
} = useForm<ApplicationFormData>({
  resolver: zodResolver(applicationSchema)
})

const onSubmit = async (data: ApplicationFormData) => {
  console.log("test")
  setServerError(null)
  try{
    await apiClient.post('/applications', data)
    router.push('/applications')
    router.refresh()
  } catch {
    setServerError('Something went wrong')
  }
}

  return (
  <div>
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="space-y-2">
        <Label htmlFor="grantId">Associated Grant</Label>
        <Controller 
          control={control}
          name="grant_id"
          render={({field}) => (
            <Combobox
             items={grants}
             itemToStringValue={(grant: Grant) => grant.title}
             onValueChange={(grant: Grant | null) => field.onChange(grant?.id ?? '')}
            >
              <ComboboxInput placeholder="Select a grant" />
              <ComboboxContent>
                <ComboboxEmpty>No grants found.</ComboboxEmpty>
                <ComboboxList>
                  {(grant) => (
                    <ComboboxItem key={grant.id} value={grant}>
                      {grant.title} - {grant.funder_name}
                    </ComboboxItem>
                  )}
                </ComboboxList>
              </ComboboxContent>
            </Combobox>
          )}>

        </Controller>
        {errors.grant_id && (
          <p className="text-sm text-red-500">{errors.grant_id.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="clientId">Associated Client</Label>
        <Controller 
          control={control}
          name="client_id"
          render={({field}) => (
            <Combobox
             items={clients}
             itemToStringValue={(client: Client) => client.name}
             onValueChange={(client: Client | null) => field.onChange(client?.id ?? '')}
            >
              <ComboboxInput placeholder="Select a client" />
              <ComboboxContent>
                <ComboboxEmpty>No clients found.</ComboboxEmpty>
                <ComboboxList>
                  {(client) => (
                    <ComboboxItem key={client.id} value={client}>
                      {client.name}
                    </ComboboxItem>
                  )}
                </ComboboxList>
              </ComboboxContent>
            </Combobox>
          )}>

        </Controller>
      </div>
      <div className="space-y-2">
        <Label htmlFor="title">Application Name</Label>
        <Input 
          id="title"
          type="text"
          placeholder="Application Name"
          {...register('title')}
        />
        {errors.title && (
          <p className="text-sm text-red-500">{errors.title.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="status">Status</Label>
        <Select onValueChange={(value) => setValue('status', value)}>
          <SelectTrigger>
            <SelectValue placeholder="Select an option" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="not_started">Not Started</SelectItem>
            <SelectItem value="draft">Draft</SelectItem>
            <SelectItem value="submitted">Submitted</SelectItem>
            <SelectItem value="approved">Approved</SelectItem>
            <SelectItem value="denied">Denied</SelectItem>
            <SelectItem value="withdrawn">Withdrawn</SelectItem>
          </SelectContent>
        </Select>
        {errors.status && (
          <p className="text-sm text-red-500">{errors.status.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="notes">Notes</Label>
        <Input 
          id="notes"
          type="text"
          placeholder="Application Notes"
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