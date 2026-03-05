'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useRouter } from 'next/navigation'
import { cn } from '@/lib/utils'
import { Users, FileText, ClipboardList, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'

const navigation = [
  { name: 'Clients', href: '/clients', icon: Users},
  { name: 'Grants', href: '/grants', icon: FileText},
  { name: 'Applications', href: '/applications', icon: ClipboardList}
]

export default function Sidebar() {
  const pathname = usePathname()
  const router = useRouter()

  const handleLogout = async () => {
    await fetch('/api/auth/logout', { method: 'POST' })
    router.push('/login')
    router.refresh()
  }

  return (
    <div className='flex flex-col h-full w-64 bg-white border-r border-gray-200'>
      <div className='flex items-center h-16 px-6 border-b border-gray-200'>
        <span className='text-lg font-semibold text-gray-900'>
          Grant Tool
        </span>
      </div>
      <nav className='flex-1 px-4 py-4 space-y-1'>
        {navigation.map((item) => {
          const isActive = pathname.startsWith(item.href)
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                'flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors',
                isActive ? 'bg-gray-100 text-gray-900' : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              )}
            >
              <item.icon className='h-5 w-5' />
              {item.name}
            </Link>
          )
        })}
      </nav>
      <div className='px-4 py-4 border-t border-gray-200'>
        <Button
          variant='ghost'
          className='w-full justify-start gap-3 text-gray-600 hover:text-gray-900'  
          onClick={handleLogout}
        >
          <LogOut className='h-5 w-5' />
          Sign Out
        </Button>
      </div>
    </div>
  )
}