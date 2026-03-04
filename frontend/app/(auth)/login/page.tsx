'use client'

import { useState } from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from 'zod'
import { useRouter } from "next/navigation"
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from '@/components/ui/card'
import Link from 'next/link'

const loginSchema = z.object({
  email: z.email('Please enter a valid email'),
  password: z.string().min(8, 'Password must be at least 8 characters')
})

type LoginFormData = z.infer<typeof loginSchema>

export default function LoginPage(){
  const router = useRouter()
  const [serverError, setServerError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginFormData) => {
    setServerError(null)
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json'},
        body: JSON.stringify(data),
      })

      const result = await response.json()

      if(!response.ok) {
        setServerError(result.error || 'Login failed')
        return
      }

      router.push('/clients')
      router.refresh()

    } catch {
      setServerError('Something went wrong. Please try again.')
    }
  }
  
  return (
  <div className="min-h-screen flex items-center justify-center bg-gray-50">
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Welcome back</CardTitle>
        <CardDescription>
          Sign in to your grant management account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="you@example.com"
              {...register('email')}
            />
            {errors.email && (
              <p className="text-sm text-red-500">{errors.email.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">Password</Label>
            <Input 
              id="password"
              type="password"
              placeholder="********"
              {...register('password')}
            />
            {errors.password && (
              <p className="text-sm text-red-500">{errors.password.message}</p>
            )}
          </div>
          {serverError && (
            <p className="text-sm text-red-500">{serverError}</p>
          )}
          <Button
            type="submit"
            className="w-full"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Signing in...' : 'Sign in'}
          </Button>
        </form>
      </CardContent>
      <CardFooter className="justify-center">
          <p className="text-sm text-gray-500">
            Don't have an account?{' '}
            <Link href="/register" className="text-blue-500 hovers:underline">
              Register
            </Link>
          </p>
      </CardFooter>
    </Card>
  </div>
)
}


