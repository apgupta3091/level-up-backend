'use client'

import { useEffect, useState } from 'react'
import { api, Progress } from '@/lib/api'
import { getAccessToken } from '@/lib/token'

export default function ProgressPage() {
  const [progress, setProgress] = useState<Progress | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = getAccessToken()
    if (!token) { setError('Not authenticated'); return }
    api.progress.get(token).then(setProgress).catch((e) => setError(e.message))
  }, [])

  if (error) return <div className="text-red-400">{error}</div>
  if (!progress) return <div className="text-gray-500">Loading...</div>

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold mb-8">Your Progress</h1>

      <div className="grid grid-cols-2 gap-4 mb-10">
        <div className="border border-gray-800 rounded-xl p-6 bg-gray-900">
          <div className="text-3xl font-bold text-indigo-400 mb-1">
            {progress.completed_lesson_ids.length}
          </div>
          <div className="text-sm text-gray-400">Lessons completed</div>
        </div>
        <div className="border border-gray-800 rounded-xl p-6 bg-gray-900">
          <div className="text-3xl font-bold text-indigo-400 mb-1">
            {progress.completed_skill_ids.length}
          </div>
          <div className="text-sm text-gray-400">Skills checked off</div>
        </div>
      </div>

      {progress.completed_lesson_ids.length === 0 && (
        <p className="text-gray-500 text-sm">
          No lessons completed yet. Start with Module 1.
        </p>
      )}
    </div>
  )
}
