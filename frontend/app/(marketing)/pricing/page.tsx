import Link from 'next/link'

export default function PricingPage() {
  return (
    <main className="min-h-screen bg-gray-950 text-white px-6 py-24">
      <div className="max-w-2xl mx-auto text-center">
        <h1 className="text-4xl font-bold mb-4">Simple pricing</h1>
        <p className="text-gray-400 mb-16">One plan. Everything included.</p>

        <div className="border border-indigo-500 rounded-2xl p-10 bg-gray-900">
          <div className="text-sm font-semibold text-indigo-400 uppercase tracking-widest mb-4">
            Founding Member
          </div>
          <div className="text-6xl font-bold mb-2">$29</div>
          <div className="text-gray-400 mb-8">per month</div>

          <ul className="text-left space-y-3 mb-10">
            {[
              'Full access to all 4 modules',
              'Weekly roadmap drops',
              'Skill tracker + progress dashboard',
              'Production-grade project assignments',
              'Rubric-based feedback on submissions',
              'Code templates & checklists',
              'Interview prep tie-in',
              'Private Discord community',
            ].map((f) => (
              <li key={f} className="flex items-center gap-3 text-sm text-gray-300">
                <span className="text-green-400 font-bold">✓</span> {f}
              </li>
            ))}
          </ul>

          <Link
            href="/sign-up"
            className="block w-full bg-indigo-500 hover:bg-indigo-400 text-white py-3 rounded-lg font-semibold text-center transition-colors"
          >
            Get started →
          </Link>
        </div>
      </div>
    </main>
  )
}
