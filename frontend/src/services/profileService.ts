import { fetchMyProfile, updateMyProfile } from "@/api/user";
import { UserProfile, UserStatistics, UserSettings } from "@/types/user";

export interface ProfileDataResponse {
  user: UserProfile;
  stats: UserStatistics;
  settings: UserSettings;
  subscriptions?: any[];
  courses?: any[];
}

export class ProfileService {
  private static instance: ProfileService;

  private constructor() {}

  public static getInstance(): ProfileService {
    if (!ProfileService.instance) {
      ProfileService.instance = new ProfileService();
    }
    return ProfileService.instance;
  }

  public async getProfileData(): Promise<ProfileDataResponse> {
    const u = await fetchMyProfile();
    
    const progressStats = u?.progress || {};
    const notifs = u?.notifications || {};
    const enrolledCourses = Array.isArray(u?.courses) ? u.courses : [];
    const subsList = Array.isArray(u?.subscriptions) ? u.subscriptions : [];
    
    const apiName = `${u?.first_name || ''} ${u?.last_name || ''}`.trim() || u?.username || u?.name || "";
    const apiTgUser = u?.username || u?.telegram_username || "";
    const apiAvatar = u?.photo_url || u?.avatar || "";

    return {
      user: {
        id: String(u?.id || u?._id || ""),
        telegramId: String(u?.telegram_id || u?.telegramId || ""),
        telegramUsername: apiTgUser,
        name: apiName,
        avatarUrl: apiAvatar,
        departmentName: u?.department_name || "",
        gradeName: u?.grade_name || "",
        joinedDate: u?.created_at ? new Date(u.created_at).toISOString().split("T")[0] : "",
        unreadNotificationsCount: Number(notifs.unread_count || u?.unreadNotificationsCount || 0),
      },
      stats: {
        enrolledCoursesCount: enrolledCourses.length || Number(u?.enrolledCoursesCount || 0),
        completedLecturesCount: Number(progressStats.completed_lessons_count || 0),
        watchTimeHours: Number(progressStats.watch_time_hours || 0),
        overallCompletionPercentage: Number(progressStats.overall_completion_percentage || 0),
      },
      settings: {
        darkMode: localStorage.getItem("wy_dark_mode") === "true",
        language: (localStorage.getItem("wy_lang") as "ar" | "en") || "ar",
        videoQuality: (localStorage.getItem("wy_video_quality") as UserSettings["videoQuality"]) || "1080p",
        autoPlayVideo: localStorage.getItem("wy_auto_play") !== "false",
        pushNotifications: localStorage.getItem("wy_notifications") !== "false",
      },
      subscriptions: subsList,
      courses: enrolledCourses,
    };
  }

  public async updateSettings(settings: Partial<UserSettings>): Promise<boolean> {
    try {
      await updateMyProfile(settings);
      return true;
    } catch (err) {
      return false;
    }
  }
}

export const profileService = ProfileService.getInstance();
