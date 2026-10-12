import { api } from "@/api";
import { Course, CourseFilterOptions } from "@/types/course";

export interface FilterMetaOption {
  id: string;
  name: string;
}

export interface CoursesFilterMetadata {
  departments: FilterMetaOption[];
  grades: FilterMetaOption[];
  subjects: FilterMetaOption[];
  instructors: FilterMetaOption[];
}

export class CourseService {
  private static instance: CourseService;

  private constructor() {}

  public static getInstance(): CourseService {
    if (!CourseService.instance) {
      CourseService.instance = new CourseService();
    }
    return CourseService.instance;
  }

  public async getFilterMetadata(): Promise<CoursesFilterMetadata> {
    try {
      const [depts, subjs, teachers] = await Promise.all([
        api.courses.fetchDepartments().catch(() => []),
        api.courses.fetchSubjects().catch(() => []),
        api.courses.fetchTeachers().catch(() => []),
      ]);

      const deptList = Array.isArray(depts) ? depts : (depts as any)?.data ?? [];
      const subjList = Array.isArray(subjs) ? subjs : (subjs as any)?.data ?? [];
      const teacherList = Array.isArray(teachers) ? teachers : (teachers as any)?.data ?? [];

      return {
        departments: deptList.map((d: any) => ({
          id: String(d.id || d._id),
          name: d.name || "القسم",
        })),
        grades: [
          { id: "grd_6", name: "السادس الإعدادي" },
          { id: "grd_5", name: "الخامس الإعدادي" },
          { id: "grd_4", name: "الرابع الإعدادي" },
        ],
        subjects: subjList.map((s: any) => ({
          id: String(s.id || s._id),
          name: s.name || "المادة",
        })),
        instructors: teacherList.map((t: any) => {
          const first = t.first_name || "";
          const last = t.last_name || "";
          const fullName = `${first} ${last}`.trim();
          return {
            id: String(t.id || t._id),
            name: fullName || t.name || t.username || "",
          };
        }),
      };
    } catch (err) {
      if (import.meta.env.DEV) console.warn("Failed to load filter metadata from API", err);
      return {
        departments: [],
        grades: [
          { id: "grd_6", name: "السادس الإعدادي" },
          { id: "grd_5", name: "الخامس الإعدادي" },
          { id: "grd_4", name: "الرابع الإعدادي" },
        ],
        subjects: [],
        instructors: [],
      };
    }
  }

  public async getCourses(filters?: Partial<CourseFilterOptions>): Promise<Course[]> {
    try {
      const res = await api.courses.fetchCoursesList({
        department_id: filters?.departmentId,
        subject_id: filters?.subjectId,
        teacher_id: filters?.instructorId,
        is_published: true,
        limit: 100,
      });

      const rawList = Array.isArray(res) ? res : (res as any)?.data ?? [];
      return rawList.map((c: any) => ({
        id: String(c.id || c._id || ""),
        title: c.title || "",
        subjectName: c.subject_name || c.subjectName || "",
        departmentName: c.department_name || c.department || "",
        gradeName: c.grade_name || c.stage || "",
        coverImageUrl: c.thumbnail_url || c.thumbnail || c.coverImageUrl || "",
        instructor: {
          id: String(c.teacher_id || ""),
          name: c.teacher_name || c.instructor_name || "",
          avatarUrl: c.teacher_avatar || "",
          title: "",
        },
        lecturesCount: Number(c.lessons_count || c.total_lessons || c.lecturesCount || 0),
        durationMinutes: Number(c.duration_minutes || 0),
        studentsCount: Number(c.students_count || 0),
        rating: Number(c.rating || 0),
        reviewsCount: Number(c.reviews_count || 0),
        price: Number(c.price || 0),
        currency: c.currency || "د.ع",
        isFree: Number(c.price || 0) === 0,
        enrollmentStatus: c.is_enrolled ? "enrolled" : "not_enrolled",
        completionPercentage: Number(c.completion_percentage || 0),
        categoryId: c.price === 0 ? "free" : "all",
        createdAt: c.created_at || "",
        viewsCount: Number(c.views_count || 0),
      }));
    } catch (err) {
      if (import.meta.env.DEV) console.error("[CourseService.getCourses] Error:", err);
      return [];
    }
  }

  public async getCourseById(courseId: string): Promise<any> {
    return await api.courses.fetchCourseById(courseId);
  }
}

export const courseService = CourseService.getInstance();
