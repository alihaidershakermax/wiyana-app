import { adminApiClient } from './adminApi'

export const isDemoMode = (): boolean => false

export const apiAdapter = {
  // -------------------------------------------------------------------------
  // Dashboard & Statistics  →  GET /dashboard/overview
  // -------------------------------------------------------------------------
  async getDashboardOverview() {
    const { data } = await adminApiClient.get('/dashboard/overview')
    return data
  },

  // -------------------------------------------------------------------------
  // Students  →  GET /students  |  GET /students/{id}  |  POST /students/{id}/suspend|activate
  // -------------------------------------------------------------------------
  async getStudents(params?: { search?: string; status?: string }) {
    const { data } = await adminApiClient.get('/students/', { params })
    return data
  },

  async getStudentDetails(id: string) {
    const { data } = await adminApiClient.get(`/students/${id}`)
    return data
  },

  /**
   * Toggle: if student is active → suspend, else → activate.
   * Backend has separate endpoints: POST /students/{id}/suspend | activate
   */
  async toggleStudentStatus(id: string, currentStatus?: string) {
    const endpoint = currentStatus?.toUpperCase() === 'ACTIVE'
      ? `/students/${id}/suspend`
      : `/students/${id}/activate`
    const { data } = await adminApiClient.post(endpoint)
    return data
  },

  // -------------------------------------------------------------------------
  // Courses  →  GET /courses  |  GET /courses/{id}  |  POST /courses
  // -------------------------------------------------------------------------
  async getCourses(params?: { search?: string; status?: string }) {
    const { data } = await adminApiClient.get('/courses/', { params })
    return data
  },

  async getCourseDetails(id: string) {
    const { data } = await adminApiClient.get(`/courses/${id}`)
    return data
  },

  async createCourse(payload: any) {
    const { data } = await adminApiClient.post('/courses/', payload)
    return data
  },

  // -------------------------------------------------------------------------
  // Subscriptions  →  GET /panel/subscriptions (admin list)
  // -------------------------------------------------------------------------
  async getSubscriptions(params?: { search?: string }) {
    const { data } = await adminApiClient.get('/panel/subscriptions', { params })
    return data
  },

  // -------------------------------------------------------------------------
  // Subscription Codes  →  GET /subscription-codes  |  POST /subscription-codes/batch-generate
  // -------------------------------------------------------------------------
  async getSubscriptionCodes(params?: { search?: string; status?: string }) {
    const { data } = await adminApiClient.get('/subscription-codes/', { params })
    return data
  },

  async batchGenerateCodes(payload: { course_ids: string[]; count: number; duration_days: number; prefix: string }) {
    const { data } = await adminApiClient.post('/subscription-codes/batch-generate', payload)
    return data
  },

  // -------------------------------------------------------------------------
  // Notifications  →  GET /notifications  |  POST /notifications/broadcast
  // -------------------------------------------------------------------------
  async getNotifications() {
    const { data } = await adminApiClient.get('/notifications/')
    return data
  },

  async broadcastNotification(payload: { title: string; message: string; target_type: string; target_id?: string }) {
    const { data } = await adminApiClient.post('/notifications/broadcast', payload)
    return data
  },

  // -------------------------------------------------------------------------
  // Analytics  →  GET /analytics/owner  |  GET /analytics/teacher
  // -------------------------------------------------------------------------
  async getOwnerAnalytics(range: string = '30d') {
    const { data } = await adminApiClient.get('/analytics/owner', { params: { range } })
    return data
  },

  async getTeacherAnalytics(range: string = '30d') {
    const { data } = await adminApiClient.get('/analytics/teacher', { params: { range } })
    return data
  },
}
