export interface Section {
  id: string
  course_id: string
  title: string
  description?: string | null
  order: number
  lessons_count: number
  created_at: string
  updated_at?: string | null
}

export interface SectionCreate {
  title: string
  description?: string
  order?: number
}

export interface SectionUpdate {
  title?: string
  description?: string
  order?: number
}

export interface SectionsReorderPayload {
  order: string[]
}
