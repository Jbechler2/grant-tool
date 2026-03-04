export interface User {
  id: string
  email: string
  role: string
}

export interface Client {
  id: string
  name: string
  contact_name: string
  contact_phone: string
  contact_email: string
  notes: string
}

export interface Grant{
  id: string
  grant_writer_id: string
  title: string
  funder_name: string
  funder_website: string
  description: string
  award_amount_min: number | null
  award_amount_max: number | null
  eligibility_notes: string
  estimated_application_hours: number | null
  visibility: string
  created_at: string
  updated_at: string
}

export interface Deadline {
  id: string
  grant_id: string
  label: string
  date: string
  description: string
  created_at: string
}

export interface Application {
  id: string
  grant_writer_id: string
  grant_id: string
  client_id: string
  title: string
  status: string
  is_exclusive: boolean
  published_at: string | null
  notes: string
  created_at: string
  updated_at: string
}

export interface AuthResponse {
    token: string
    email: string
    role: string
}

export interface ApiError {
    error: string
}