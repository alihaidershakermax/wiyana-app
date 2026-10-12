export interface Teacher {
  id: string
  telegram_id?: number | null
  username?: string | null
  first_name?: string | null
  last_name?: string | null
  bio?: string | null
  role: string
  is_active: boolean
  created_at?: string | null
}

export interface TeacherCreate {
  name?: string
  username?: string
  telegram_id?: number | null
  bio?: string
}
