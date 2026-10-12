export type SubjectStatus = 'active' | 'inactive' | 'archived'

export interface Subject {
  id: string
  name: string
  code: string
  description?: string | null
  department?: string | null
  department_id?: string | null
  stage?: string | null
  status: SubjectStatus
  created_by?: string | null
  courses_count: number
  students_count: number
  created_at: string
  updated_at?: string | null
}

export interface SubjectCreate {
  name: string
  code: string
  description?: string
  department?: string
  stage?: string
  status?: SubjectStatus
}

export interface SubjectUpdate {
  name?: string
  code?: string
  description?: string
  department?: string
  stage?: string
  status?: SubjectStatus
}

export interface SubjectListFilters {
  search?: string
  department?: string
  stage?: string
  status?: SubjectStatus
  page?: number
  limit?: number
}

export interface DeleteSubjectResult {
  deleted: boolean
  archived: boolean
  message: string
  courses_count?: number
}
