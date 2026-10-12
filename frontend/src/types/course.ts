export interface Instructor {
  id: string;
  name: string;
  avatarUrl?: string;
  title?: string;
}

export type CourseEnrollmentStatus = "not_enrolled" | "enrolled" | "completed";

export interface Course {
  id: string;
  title: string;
  subjectName: string;
  departmentName?: string; // e.g., "الفرع العلمي"
  gradeName?: string; // e.g., "السادس الاعدادي"
  coverImageUrl: string;
  instructor: Instructor;
  lecturesCount: number;
  durationMinutes: number;
  studentsCount: number;
  rating: number;
  reviewsCount: number;
  price: number;
  currency: string;
  isFree: boolean;
  enrollmentStatus: CourseEnrollmentStatus;
  completionPercentage?: number;
  lastWatchedLectureId?: string;
  lastWatchedLectureTitle?: string;
  categoryId: string;
  createdAt: string;
  viewsCount: number;
}

export type CourseTabCategory = "all" | "my_courses" | "newest" | "free" | "paid" | "popular" | "favorites";

export type CategoryFilterId = "all" | "new" | "free" | "paid" | "popular";

export type CourseSortOption = "newest" | "oldest" | "popular" | "top_rated" | "alphabetical";

export interface CourseFilterOptions {
  searchQuery: string;
  categoryTab: CourseTabCategory;
  departmentId?: string;
  gradeId?: string;
  subjectId?: string;
  instructorId?: string;
  isFreeOnly?: boolean;
  isPaidOnly?: boolean;
  isSubscribedOnly?: boolean;
  sortBy: CourseSortOption;
}

export interface CategoryOption {
  id: CourseTabCategory;
  label: string;
  iconName?: string;
}
