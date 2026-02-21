'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { api, ModuleDetail } from '@/lib/api'
import { getAccessToken } from '@/lib/token'

export default function ModulePage() {
  const { slug } = useParams<{ slug: string }>()
  const [module, setModule] = useState<ModuleDetail | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = getAccessToken()
    if (!token) { setError('Not authenticated'); return }
    api.modules.get(slug, token)
      .then(setModule)
      .catch((e) => setError(e.message))
  }, [slug])

  if (error) return <div className="text-red-400">{error}</div>
  if (!module) return <div className="text-gray-500">Loading...</div>

  return (
    <div className="max-w-3xl">
      <div className="text-xs font-bold text-indigo-400 uppercase tracking-widest mb-2">Module</div>
      <h1 className="text-2xl font-bold mb-2">{module.title}</h1>
      <p className="text-gray-400 mb-2">{module.description}</p>
      <div className="text-sm text-gray-500 mb-8">
        {module.completed_lessons}/{module.total_lessons} lessons complete ·{' '}
        {module.estimated_hours}h estimated
      </div>

      {/* Progress bar */}
      <div className="w-full h-1.5 bg-gray-800 rounded-full mb-10">
        <div
          className="h-1.5 bg-indigo-500 rounded-full transition-all"
          style={{ width: `${module.total_lessons ? (module.completed_lessons / module.total_lessons) * 100 : 0}%` }}
        />
      </div>

      {/* Lesson list */}
      <div className="space-y-2 mb-10">
        {module.lessons.map((lesson, i) => (
          <Link
            key={lesson.id}
            href={`/modules/${slug}/lessons/${lesson.slug}`}
            className="flex items-center justify-between border border-gray-800 rounded-lg px-4 py-3 hover:border-gray-600 hover:bg-gray-900 transition-colors"
          >
            <div className="flex items-center gap-3">
              <span className="text-xs text-gray-600 w-5 text-right">{i + 1}</span>
              <span className="text-sm font-medium">{lesson.title}</span>
            </div>
            <span className="text-xs text-gray-500">{lesson.estimated_minutes}m</span>
          </Link>
        ))}
      </div>

      {/* Assignment link */}
      <Link
        href={`/modules/${slug}/assignment`}
        className="inline-flex items-center gap-2 border border-indigo-700 text-indigo-400 hover:bg-indigo-900/30 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
      >
        View assignment →
      </Link>
    </div>
  )
}
