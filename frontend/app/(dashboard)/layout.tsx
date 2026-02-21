import Link from 'next/link'

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-950 text-white flex">
      {/* Sidebar */}
      <aside className="w-56 border-r border-gray-800 flex flex-col p-4 gap-1 shrink-0">
        <div className="font-bold text-sm tracking-tight px-3 py-2 mb-4 text-gray-300">
          Level Up Backend
        </div>
        <NavLink href="/dashboard">Dashboard</NavLink>
        <NavLink href="/modules/go-concurrency">Module 1</NavLink>
        <NavLink href="/progress">Progress</NavLink>
        <NavLink href="/submissions">Submissions</NavLink>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-y-auto p-8">
        {children}
      </main>
    </div>
  )
}

function NavLink({ href, children }: { href: string; children: React.ReactNode }) {
  return (
    <Link
      href={href}
      className="text-sm text-gray-400 hover:text-white hover:bg-gray-800 px-3 py-2 rounded-md transition-colors"
    >
      {children}
    </Link>
  )
}
