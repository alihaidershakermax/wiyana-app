import { api } from "@/api";
import { HomeDashboardData } from "@/types/home";
import { Course } from "@/types/course";

function mapRawCourse(c: any): Course {
  return {
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
    enrollmentStatus: c.enrollmentStatus || (c.is_enrolled ? "enrolled" : "not_enrolled"),
    completionPercentage: Number(c.completionPercentage || c.completion_percentage || 0),
    lastWatchedLectureTitle: c.lastWatchedLectureTitle,
    categoryId: c.price === 0 ? "free" : "all",
    createdAt: c.created_at || "",
    viewsCount: Number(c.views_count || 0),
  };
}

export class HomeService {
  private static instance: HomeService;

  private constructor() {}

  public static getInstance(): HomeService {
    if (!HomeService.instance) {
      HomeService.instance = new HomeService();
    }
    return HomeService.instance;
  }

  public async getHomeDashboard(): Promise<HomeDashboardData> {
    const data = await api.home.fetchHomeDashboard();
    
    const u = data?.user || {};
    const announcement = data?.announcement || null;

    const continueList = Array.isArray(data?.continueLearningCourses)
      ? data.continueLearningCourses.map(mapRawCourse)
      : [];
    const recommendedList = Array.isArray(data?.recommendedCourses)
      ? data.recommendedCourses.map(mapRawCourse)
      : [];
    const latestList = Array.isArray(data?.latestCourses)
      ? data.latestCourses.map(mapRawCourse)
      : [];
    const popularList = Array.isArray(data?.popularCourses)
      ? data.popularCourses.map(mapRawCourse)
      : [];

    return {
      user: {
        id: String(u.id || ""),
        telegramId: String(u.telegramId || u.telegram_id || ""),
        telegramUsername: u.username || "",
        name: u.name || "",
        avatarUrl: u.avatarUrl || u.photo_url || "",
        departmentName: u.departmentName || "",
        gradeName: u.gradeName || "",
        joinedDate: u.created_at ? String(u.created_at).slice(0, 10) : "",
        enrolledCoursesCount: Number(u.enrolledCoursesCount || 0),
        overallCompletionPercentage: Number(u.overallCompletionPercentage || 0),
        unreadNotificationsCount: Number(u.unreadNotificationsCount || 0),
      },
      announcement,
      continueLearningCourses: continueList,
      recommendedCourses: recommendedList,
      latestCourses: latestList,
      popularCourses: popularList,
    };
  }
}

export const homeService = HomeService.getInstance();
