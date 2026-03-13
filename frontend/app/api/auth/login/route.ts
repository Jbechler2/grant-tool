import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  try {
    const body = await request.json()

    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
        credentials: 'include'
      }
    )

    const data = await response.json()
    if(!response.ok) {
      return NextResponse.json(data, { status: response.status })
    }

    const nextResponse = NextResponse.json({
      email: data.email,
      role: data.role,
    })

    const cookieHeaders = response.headers.getSetCookie()
    cookieHeaders.forEach(cookie => {
      nextResponse.headers.append('set-cookie', cookie)
    });

    return nextResponse
  } catch {
    return NextResponse.json(
      { error: 'Login failed' },
      { status: 500 }
    )
  }
}