export interface UserProfile {
  id: string;
  telegramId: string;
  telegramUsername?: string;
  name: string;
  avatarUrl?: string;
  departmentName: string;
  gradeName: string;
  joinedDate: string;
  unreadNotificationsCount: number;
  enrolledCoursesCount?: number;
  overallCompletionPercentage?: number;
}

export interface UserStatistics {
  enrolledCoursesCount: number;
  completedLecturesCount: number;
  watchTimeHours: number;
  overallCompletionPercentage: number;
}

export interface UserSettings {
  darkMode: boolean;
  language: "ar" | "en";
  videoQuality: "auto" | "1080p" | "720p" | "480p" | "360p";
  autoPlayVideo: boolean;
  pushNotifications: boolean;
}
