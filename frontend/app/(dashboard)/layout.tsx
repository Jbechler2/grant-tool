import Sidebar from '@/components/features/navigation/sidebar'

export default function DashBoardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
     <div className='flex h-screen bg-gray-50'>
        <Sidebar />
        <main className='flex-1 overflow-y-auto'>
          <div className="p-8">
            {children}
          </div>
        </main>
     </div>
  )
}