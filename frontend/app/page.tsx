import Link from 'next/link'

export default function HomePage() {
  return (
    <main className="min-h-screen bg-gray-950 text-white">
      {/* Nav */}
      <nav className="border-b border-gray-800 px-6 py-4 flex items-center justify-between max-w-6xl mx-auto">
        <span className="font-bold text-lg tracking-tight">Level Up Backend</span>
        <div className="flex items-center gap-4">
          <Link href="/pricing" className="text-sm text-gray-400 hover:text-white transition-colors">
            Pricing
          </Link>
          <Link href="/sign-in" className="text-sm text-gray-400 hover:text-white transition-colors">
            Sign in
          </Link>
          <Link
            href="/sign-up"
            className="text-sm bg-white text-gray-950 px-4 py-1.5 rounded-md font-medium hover:bg-gray-100 transition-colors"
          >
            Get started
          </Link>
        </div>
      </nav>

      {/* Hero */}
      <section className="max-w-4xl mx-auto px-6 pt-24 pb-16 text-center">
        <div className="inline-block text-xs font-semibold tracking-widest text-indigo-400 uppercase mb-6 border border-indigo-800 px-3 py-1 rounded-full">
          Backend Engineering
        </div>
        <h1 className="text-5xl sm:text-6xl font-bold leading-tight mb-6">
          From mid-level to<br />
          <span className="text-indigo-400">senior backend</span><br />
          in 6 months.
        </h1>
        <p className="text-xl text-gray-400 max-w-2xl mx-auto mb-10">
          A structured, production-grade roadmap for backend engineers who want FAANG-tier skills.
          Not theory. Execution.
        </p>
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Link
            href="/sign-up"
            className="bg-indigo-500 hover:bg-indigo-400 text-white px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
          >
            Start your roadmap →
          </Link>
          <Link
            href="/pricing"
            className="text-gray-400 hover:text-white px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
          >
            See pricing
          </Link>
        </div>
      </section>

      {/* Phases */}
      <section className="max-w-5xl mx-auto px-6 py-16">
        <h2 className="text-2xl font-bold text-center mb-12 text-gray-300">The Skill Ladder</h2>
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {[
            {
              phase: '01',
              title: 'Concurrency',
              items: ['Worker pools', 'Context propagation', 'Backpressure', 'Goroutine leaks', 'Graceful shutdown'],
            },
            {
              phase: '02',
              title: 'Distributed Systems',
              items: ['Idempotency', 'Retries', 'Circuit breakers', 'Rate limiting', 'Message queues'],
            },
            {
              phase: '03',
              title: 'Reliability',
              items: ['Observability', 'Structured logging', 'Tracing', 'SLOs', 'Deploy strategies'],
            },
            {
              phase: '04',
              title: 'Architecture',
              items: ['Tradeoff simulations', 'CAP theorem', 'Data modeling', 'Scaling decisions'],
            },
          ].map((p) => (
            <div key={p.phase} className="border border-gray-800 rounded-xl p-6 bg-gray-900">
              <div className="text-xs font-bold text-indigo-400 mb-2">PHASE {p.phase}</div>
              <h3 className="font-bold text-lg mb-4">{p.title}</h3>
              <ul className="space-y-1.5">
                {p.items.map((item) => (
                  <li key={item} className="text-sm text-gray-400 flex items-center gap-2">
                    <span className="text-indigo-500">→</span> {item}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="max-w-2xl mx-auto px-6 py-16 text-center">
        <h2 className="text-3xl font-bold mb-4">Ready to level up?</h2>
        <p className="text-gray-400 mb-8">Join engineers building real systems, not just studying theory.</p>
        <Link
          href="/sign-up"
          className="bg-indigo-500 hover:bg-indigo-400 text-white px-8 py-3 rounded-lg font-semibold text-lg transition-colors"
        >
          Get started →
        </Link>
      </section>

      <footer className="border-t border-gray-800 py-8 text-center text-sm text-gray-600">
        © {new Date().getFullYear()} Level Up Backend
      </footer>
    </main>
  )
}
