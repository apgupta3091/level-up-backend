const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
  }
}

async function request<T>(path: string, options: RequestInit = {}, token?: string): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const res = await fetch(`${API_URL}${path}`, { ...options, headers })
  const data = await res.json()

  if (!res.ok) {
    throw new ApiError(res.status, data.error ?? 'Something went wrong')
  }
  return data as T
}

// ── Auth ──────────────────────────────────────────────────────

export type AuthResponse = {
  access_token: string
  refresh_token: string
  user: User
}

export type User = {
  id: string
  email: string
  name: string
  subscription_status: 'free' | 'active' | 'cancelled' | 'past_due'
}

export const api = {
  auth: {
    register: (email: string, password: string, name: string) =>
      request<AuthResponse>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password, name }),
      }),

    login: (email: string, password: string) =>
      request<AuthResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      }),

    refresh: (refreshToken: string) =>
      request<{ access_token: string; refresh_token: string }>('/auth/refresh', {
        method: 'POST',
        body: JSON.stringify({ refresh_token: refreshToken }),
      }),
  },

  // ── Modules ────────────────────────────────────────────────

  modules: {
    list: (token: string) =>
      request<{ modules: Module[] }>('/modules', {}, token),

    get: (slug: string, token: string) =>
      request<ModuleDetail>(`/modules/${slug}`, {}, token),
  },

  // ── Lessons ────────────────────────────────────────────────

  lessons: {
    get: (moduleSlug: string, lessonSlug: string, token: string) =>
      request<Lesson>(`/modules/${moduleSlug}/lessons/${lessonSlug}`, {}, token),

    complete: (lessonId: string, token: string) =>
      request<{ status: string }>(`/lessons/${lessonId}/complete`, { method: 'POST' }, token),
  },

  // ── Progress ───────────────────────────────────────────────

  progress: {
    get: (token: string) =>
      request<Progress>('/progress', {}, token),
  },

  // ── Skills ─────────────────────────────────────────────────

  skills: {
    list: (moduleSlug: string, token: string) =>
      request<{ skills: Skill[] }>(`/modules/${moduleSlug}/skills`, {}, token),

    complete: (skillId: string, token: string) =>
      request<{ status: string }>(`/skills/${skillId}/complete`, { method: 'POST' }, token),
  },

  // ── Submissions ────────────────────────────────────────────

  submissions: {
    list: (token: string) =>
      request<{ submissions: Submission[] }>('/submissions', {}, token),

    get: (id: string, token: string) =>
      request<Submission>(`/submissions/${id}`, {}, token),

    create: (assignmentId: string, githubUrl: string, writtenAnswers: string, token: string) =>
      request<Submission>('/submissions', {
        method: 'POST',
        body: JSON.stringify({ assignment_id: assignmentId, github_url: githubUrl, written_answers: writtenAnswers }),
      }, token),
  },

  // ── Payments ───────────────────────────────────────────────

  payments: {
    checkout: (token: string) =>
      request<{ url: string }>('/payments/checkout', { method: 'POST' }, token),

    subscription: (token: string) =>
      request<{ subscription_status: string; stripe_subscription_id: string | null }>(
        '/payments/subscription', {}, token
      ),
  },
}

// ── Types ─────────────────────────────────────────────────────

export type Module = {
  id: string
  title: string
  slug: string
  description: string
  order_index: number
  estimated_hours: number
}

export type ModuleDetail = Module & {
  total_lessons: number
  completed_lessons: number
  lessons: LessonSummary[]
}

export type LessonSummary = {
  id: string
  title: string
  slug: string
  order_index: number
  estimated_minutes: number
}

export type Lesson = {
  id: string
  module_id: string
  title: string
  slug: string
  content: string
  order_index: number
  estimated_minutes: number
}

export type Skill = {
  id: string
  skill_name: string
  order_index: number
}

export type Progress = {
  completed_lesson_ids: string[]
  completed_skill_ids: string[]
}

export type Submission = {
  id: string
  assignment_id: string
  user_id: string
  github_url: string
  written_answers: string
  status: 'pending' | 'reviewed' | 'approved' | 'needs_revision'
  feedback: string | null
  submitted_at: string
  reviewed_at: string | null
}

export { ApiError }
