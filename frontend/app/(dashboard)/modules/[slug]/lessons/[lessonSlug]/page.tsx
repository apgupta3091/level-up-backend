'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import { api, Lesson } from '@/lib/api'
import { getAccessToken } from '@/lib/token'

export default function LessonPage() {
  const { slug, lessonSlug } = useParams<{ slug: string; lessonSlug: string }>()
  const [lesson, setLesson] = useState<Lesson | null>(null)
  const [completed, setCompleted] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const token = getAccessToken()
    if (!token) { setError('Not authenticated'); return }
    api.lessons.get(slug, lessonSlug, token)
      .then(setLesson)
      .catch((e) => setError(e.message))
  }, [slug, lessonSlug])

  async function handleComplete() {
    const token = getAccessToken()
    if (!token || !lesson) return
    await api.lessons.complete(lesson.id, token)
    setCompleted(true)
  }

  if (error) return <div className="text-red-400">{error}</div>
  if (!lesson) return <div className="text-gray-500">Loading...</div>

  return (
    <div className="max-w-2xl">
      <Link href={`/modules/${slug}`} className="text-sm text-gray-500 hover:text-gray-300 mb-6 inline-block">
        ← Back to module
      </Link>

      <h1 className="text-2xl font-bold mb-1">{lesson.title}</h1>
      <div className="text-sm text-gray-500 mb-8">{lesson.estimated_minutes} min read</div>

      {/* Content — rendered as preformatted for now, swap in a markdown renderer */}
      <div className="prose prose-invert prose-sm max-w-none">
        <pre className="whitespace-pre-wrap font-sans text-gray-300 leading-relaxed">
          {lesson.content}
        </pre>
      </div>

      <div className="mt-10 pt-8 border-t border-gray-800">
        {completed ? (
          <div className="flex items-center gap-2 text-green-400 font-medium">
            <span>✓</span> Lesson marked complete
          </div>
        ) : (
          <button
            onClick={handleComplete}
            className="bg-indigo-500 hover:bg-indigo-400 text-white px-6 py-2 rounded-lg font-medium transition-colors"
          >
            Mark as complete
          </button>
        )}
      </div>
    </div>
  )
}
