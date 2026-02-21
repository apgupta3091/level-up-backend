import Link from 'next/link'

export default function DashboardPage() {
  return (
    <div className="max-w-4xl">
      <h1 className="text-2xl font-bold mb-1">Welcome back.</h1>
      <p className="text-gray-400 mb-10">Here&apos;s where you left off.</p>

      <div className="grid sm:grid-cols-2 gap-4 mb-10">
        {[
          { phase: '01', title: 'Go Concurrency & Graceful Shutdown', slug: 'go-concurrency', available: true },
          { phase: '02', title: 'Distributed Systems Fundamentals', slug: 'distributed-systems', available: false },
          { phase: '03', title: 'Reliability & Observability', slug: 'reliability', available: false },
          { phase: '04', title: 'Architecture & Systems Thinking', slug: 'architecture', available: false },
        ].map((m) => (
          <div
            key={m.slug}
            className={`border rounded-xl p-6 ${m.available ? 'border-gray-700 bg-gray-900 hover:border-indigo-600 transition-colors' : 'border-gray-800 bg-gray-900/50 opacity-60'}`}
          >
            <div className="text-xs font-bold text-indigo-400 mb-2">PHASE {m.phase}</div>
            <h3 className="font-semibold mb-4">{m.title}</h3>
            {m.available ? (
              <Link href={`/modules/${m.slug}`} className="text-sm text-indigo-400 hover:text-indigo-300 font-medium">
                Open module →
              </Link>
            ) : (
              <span className="text-xs text-gray-600 uppercase tracking-wider">Coming soon</span>
            )}
          </div>
        ))}
      </div>

      <div className="flex gap-4">
        <Link href="/progress" className="text-sm border border-gray-700 rounded-lg px-4 py-2 text-gray-400 hover:text-white hover:border-gray-600 transition-colors">
          View progress
        </Link>
        <Link href="/submissions" className="text-sm border border-gray-700 rounded-lg px-4 py-2 text-gray-400 hover:text-white hover:border-gray-600 transition-colors">
          My submissions
        </Link>
      </div>
    </div>
  )
}
