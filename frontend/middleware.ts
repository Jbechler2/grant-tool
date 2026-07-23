import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ['/login', '/register']

export async function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value
  const refreshToken = request.cookies.get('refresh_token')?.value
  const pathname = request.nextUrl.pathname
  const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route))

  if(!token && refreshToken && !isPublicRoute) {
    const refreshRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: { Cookie: `refresh_token=${refreshToken}`}
    })

    if (refreshRes.ok) {
      const setCookies = refreshRes.headers.getSetCookie()

      const pairs = setCookies.map(c => c.split(';')[0].split('='))
      const newToken   = pairs.find(([n]) => n === 'token')?.[1]
      const newRefresh = pairs.find(([n]) => n === 'refresh_token')?.[1]

        if (newToken)   request.cookies.set('token', newToken)
        if (newRefresh) request.cookies.set('refresh_token', newRefresh)

      const response = NextResponse.next({
        request: { headers: request.headers },
      })


      for (const c of setCookies) response.headers.append('set-cookie', c)

      return response
    }

    return NextResponse.redirect(new URL('/login', request.url))
  }

  if(!token && !refreshToken && !isPublicRoute) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  if(token && isPublicRoute) {
    return NextResponse.redirect(new URL('/clients', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}