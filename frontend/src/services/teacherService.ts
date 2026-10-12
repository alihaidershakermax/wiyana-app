import { apiClient } from "@/api/client";
import { TeacherAccount, StudentQuestion } from "@/types/teacher";

interface RawQuestion {
  id: string;
  lesson_id?: string;
  user_id?: string;
  course_id?: string;
  question_text?: string;
  answer_text?: string | null;
  answered_by?: string | null;
  status?: string;
  created_at?: string;
  updated_at?: string;
}

export const teacherService = {
  // 1. Fetch all teachers (For Platform Owner)
  async getAllTeachers(): Promise<TeacherAccount[]> {
    const response = await apiClient.get<TeacherAccount[]>("/teachers/");
    return response.data;
  },

  // 2. Create New Teacher Account
  async createTeacher(payload: {
    name: string;
    telegramId: string;
    telegramUsername?: string;
    specialty: string;
    phone?: string;
    email?: string;
    commissionPct: number;
    assignedCourses: string[];
  }): Promise<TeacherAccount> {
    const response = await apiClient.post<TeacherAccount>("/users/teacher", payload);
    return response.data;
  },

  // 3. Toggle Teacher Status (Active/Suspended)
  async toggleTeacherStatus(teacherId: string, isActive: boolean): Promise<boolean> {
    await apiClient.patch(`/teachers/${teacherId}/status`, { is_active: isActive });
    return true;
  },

  // 4. Student Questions for Teacher
  async getTeacherQuestions(): Promise<StudentQuestion[]> {
    const response = await apiClient.get<RawQuestion[]>("/questions/teacher");
    const items = Array.isArray(response.data) ? response.data : [];

    const enriched = await Promise.all(
      items.map(async (q) => {
        let lessonTitle = "";
        if (q.lesson_id) {
          try {
            const lessonRes = await apiClient.get<{ title?: string }>(`/lessons/${q.lesson_id}`);
            lessonTitle = lessonRes.data?.title || "";
          } catch {
            // lesson may be inaccessible; keep empty title
          }
        }

        return {
          id: q.id,
          studentId: q.user_id || "",
          studentName: `طالب #${(q.user_id || "؟").slice(-6)}`,
          courseId: q.course_id || "",
          courseTitle: "",
          lessonId: q.lesson_id || "",
          lessonTitle,
          timestampInVideo: "",
          questionText: q.question_text || "",
          replyText: q.answer_text || undefined,
          repliedAt: q.answered_by && q.updated_at
            ? new Date(q.updated_at).toLocaleString("ar")
            : undefined,
          status: q.status === "answered" ? "answered" : "open",
          createdAt: q.created_at
            ? new Date(q.created_at).toLocaleString("ar")
            : "",
        } as StudentQuestion;
      })
    );

    return enriched;
  },

  // 5. Reply to student question
  async replyToQuestion(questionId: string, replyText: string): Promise<boolean> {
    await apiClient.post(`/questions/${questionId}/reply`, { replyText });
    return true;
  },
};
