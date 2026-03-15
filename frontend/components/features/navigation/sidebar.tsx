'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useRouter } from 'next/navigation'
import { cn } from '@/lib/utils'
import { Users, FileText, ClipboardList, LogOut, ChevronLeft, ChevronRight, Menu, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import apiClient from '@/lib/api'
import { useState } from 'react'

const navigation = [
  { name: 'Clients', href: '/clients', icon: Users},
  { name: 'Grants', href: '/grants', icon: FileText},
  { name: 'Applications', href: '/applications', icon: ClipboardList}
]

export default function Sidebar() {
  const pathname = usePathname()
  const router = useRouter()
  const [isCollapsed, setIsCollapsed] = useState(false)
  const [mobileOpen, setMobileOpen] = useState(false)

  const handleLogout = async () => {
    await apiClient.post('/auth/logout')
    router.push('/login')
    router.refresh()
  }

  return (
    <>
      {/* Mobile hamburger */}
      <button
        onClick={() => setMobileOpen(true)}
        className='md:hidden fixed top-4 left-4 z-40 p-2 rounded-md bg-white border border-gray-200 text-gray-600 shadow-sm'
      >
        <Menu className='h-5 w-5' />
      </button>

      {/* Mobile backdrop */}
      {mobileOpen && (
        <div
          className='md:hidden fixed inset-0 z-40 bg-black/30'
          onClick={() => setMobileOpen(false)}
        />
      )}

      {/* Sidebar */}
      <div className={cn(
        'flex flex-col bg-white border-r border-gray-200 transition-all duration-300',
        'fixed inset-y-0 left-0 z-50 h-full md:relative md:z-auto md:inset-auto',
        mobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0',
        isCollapsed ? 'md:w-16 w-64' : 'w-64'
      )}>
        <div className='flex items-center h-16 px-4 border-b border-gray-200'>
          {!isCollapsed && (
            <span className='flex-1 text-lg font-semibold text-gray-900'>Grant Tool</span>
          )}
          {/* Desktop collapse toggle */}
          <button
            onClick={() => setIsCollapsed(!isCollapsed)}
            className='hidden md:flex items-center justify-center p-1 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100 ml-auto'
          >
            {isCollapsed ? <ChevronRight className='h-5 w-5' /> : <ChevronLeft className='h-5 w-5' />}
          </button>
          {/* Mobile close button */}
          <button
            onClick={() => setMobileOpen(false)}
            className='md:hidden flex items-center justify-center p-1 ml-auto rounded-md text-gray-400 hover:text-gray-600'
          >
            <X className='h-5 w-5' />
          </button>
        </div>

        <nav className='flex-1 px-2 py-4 space-y-1'>
          {navigation.map((item) => {
            const isActive = pathname.startsWith(item.href)
            return (
              <div className='flex flex-row items-center justify-between' key={item.name}>
                <Link
                  href={item.href}
                  onClick={() => setMobileOpen(false)}
                  className={cn(
                    'flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors flex-1 min-w-0',
                    isActive ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900',
                    isCollapsed && 'md:justify-center'
                  )}
                  title={isCollapsed ? item.name : undefined}
                >
                  <item.icon className='h-5 w-5 flex-shrink-0' />
                  <span className={cn(isCollapsed && 'md:hidden')}>{item.name}</span>
                </Link>
                <Link
                  href={item.href + '/new'}
                  className={cn('bg-gray-200 border border-gray-400 px-2 py-1 rounded-md ml-1', isCollapsed && 'md:hidden')}
                >
                  <div className='text-lg'>+</div>
                </Link>
              </div>
            )
          })}
        </nav>

        <div className='px-2 py-4 border-t border-gray-200'>
          <Button
            variant='ghost'
            className={cn(
              'w-full gap-3 text-gray-600 hover:text-gray-900',
              isCollapsed ? 'md:justify-center md:px-2' : 'justify-start'
            )}
            onClick={handleLogout}
          >
            <LogOut className='h-5 w-5 flex-shrink-0' />
            <span className={cn(isCollapsed && 'md:hidden')}>Sign Out</span>
          </Button>
        </div>
      </div>
    </>
  )
}
