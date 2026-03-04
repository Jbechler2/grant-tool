import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    const response = await fetch (
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      }
    )

    const data = await response.json()

    if(!response.ok){
      return NextResponse.json(data, { status: response.status })
    }

    const nextResponse = NextResponse.json({
      email: data.email,
      role: data.role,
    })

    nextResponse.cookies.set('token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 15,
      path: '/'
    })

    return nextResponse
  } catch {
    return NextResponse.json(
      {error: 'Login failed'},
      { status: 500 }
    )
  }
}