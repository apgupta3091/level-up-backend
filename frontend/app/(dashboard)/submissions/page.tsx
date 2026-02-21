'use client'

import { useEffect, useState } from 'react'
import { api, Submission } from '@/lib/api'
import { getAccessToken } from '@/lib/token'

const statusColors: Record<Submission['status'], string> = {
  pending:        'text-yellow-400 bg-yellow-900/30 border-yellow-800',
  reviewed:       'text-blue-400 bg-blue-900/30 border-blue-800',
  approved:       'text-green-400 bg-green-900/30 border-green-800',
  needs_revision: 'text-red-400 bg-red-900/30 border-red-800',
}

export default function SubmissionsPage() {
  const [submissions, setSubmissions] = useState<Submission[]>([])
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = getAccessToken()
    if (!token) { setError('Not authenticated'); return }
    api.submissions.list(token)
      .then((d) => setSubmissions(d.submissions ?? []))
      .catch((e) => setError(e.message))
  }, [])

  if (error) return <div className="text-red-400">{error}</div>

  return (
    <div className="max-w-3xl">
      <h1 className="text-2xl font-bold mb-8">My Submissions</h1>

      {submissions.length === 0 ? (
        <p className="text-gray-500 text-sm">No submissions yet. Complete a module assignment to submit.</p>
      ) : (
        <div className="space-y-4">
          {submissions.map((s) => (
            <div key={s.id} className="border border-gray-800 rounded-xl p-5 bg-gray-900">
              <div className="flex items-center justify-between mb-3">
                <a
                  href={s.github_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-indigo-400 hover:underline font-medium"
                >
                  {s.github_url}
                </a>
                <span className={`text-xs px-2 py-0.5 rounded border font-medium ${statusColors[s.status]}`}>
                  {s.status.replace('_', ' ')}
                </span>
              </div>
              {s.feedback && (
                <div className="text-sm text-gray-300 bg-gray-800 rounded-lg p-3 mt-3">
                  <div className="text-xs text-gray-500 mb-1 uppercase tracking-wider">Feedback</div>
                  {s.feedback}
                </div>
              )}
              <div className="text-xs text-gray-600 mt-3">
                Submitted {new Date(s.submitted_at).toLocaleDateString()}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
