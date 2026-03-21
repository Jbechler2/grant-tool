import { refresh } from "next/cache";
import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ['/login', '/register']

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token')?.value
  const refreshToken = request.cookies.get('refresh_token')
  const pathname = request.nextUrl.pathname
  const isPublicRoute = publicRoutes.some(route => pathname.startsWith(route))

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