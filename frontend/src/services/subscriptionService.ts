import { fetchMySubscriptions, redeemSubscriptionCode } from "@/api/subscription";
import { fetchMyProfile } from "@/api/user";
import { SubscriptionInfo, CourseSubscriptionItem } from "@/types/subscription";

export class SubscriptionService {
  private static instance: SubscriptionService;

  private constructor() {}

  public static getInstance(): SubscriptionService {
    if (!SubscriptionService.instance) {
      SubscriptionService.instance = new SubscriptionService();
    }
    return SubscriptionService.instance;
  }

  public parseSubscriptionList(subs: any[]): SubscriptionInfo | null {
    if (!Array.isArray(subs) || subs.length === 0) {
      return null;
    }

    const parsedList: CourseSubscriptionItem[] = subs.map((sub: any, idx: number) => {
      const sDate = sub.start_date || sub.startDate || sub.created_at || sub.started_at;
      const eDate = sub.end_date || sub.endDate || sub.expires_at;
      
      let remDays = sub.remaining_days ?? sub.remainingDays;
      let totDays = sub.total_days ?? sub.totalDays;
      
      if ((remDays === undefined || remDays === null) && eDate) {
        const end = new Date(eDate);
        const now = new Date();
        remDays = Math.max(0, Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)));
      }
      
      if ((totDays === undefined || totDays === null) && sDate && eDate) {
        const start = new Date(sDate);
        const end = new Date(eDate);
        totDays = Math.max(1, Math.round((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)));
      }

      const startObj = sDate ? new Date(sDate) : null;
      const endObj = eDate ? new Date(eDate) : null;

      const statusStr = String(sub.status || "").toLowerCase();
      const isExp = statusStr === "expired" || (endObj ? new Date() > endObj : false);

      const title = sub.course?.title || sub.course_title || sub.plan_name || sub.planName || sub.title || "مساق دراسي";

      return {
        id: String(sub.id || sub._id || `sub_${idx}`),
        courseTitle: title,
        startDate: startObj && !isNaN(startObj.getTime()) ? startObj.toISOString().split("T")[0] : "",
        endDate: endObj && !isNaN(endObj.getTime()) ? endObj.toISOString().split("T")[0] : "",
        remainingDays: Number(remDays) || 0,
        totalDays: Number(totDays) || 30,
        status: isExp ? "expired" : "active",
      };
    });

    if (parsedList.length === 0) return null;

    // Pick active subscription first, else first item
    const activeSub = parsedList.find((s) => s.status === "active") || parsedList[0];

    return {
      id: activeSub.id,
      status: activeSub.status,
      planName: activeSub.courseTitle,
      startDate: activeSub.startDate,
      endDate: activeSub.endDate,
      remainingDays: activeSub.remainingDays,
      totalDays: activeSub.totalDays,
      subscriptionsList: parsedList,
    };
  }

  public async getSubscriptionInfo(): Promise<SubscriptionInfo | null> {
    try {
      const subs = await fetchMySubscriptions();
      if (Array.isArray(subs) && subs.length > 0) {
        const parsed = this.parseSubscriptionList(subs);
        if (parsed) return parsed;
      }
    } catch (err) {
      console.warn("fetchMySubscriptions failed, attempting profile fallback:", err);
    }

    // Check user profile for enrolled courses / subscriptions
    try {
      const profile = await fetchMyProfile();
      if (Array.isArray(profile?.subscriptions) && profile.subscriptions.length > 0) {
        const parsed = this.parseSubscriptionList(profile.subscriptions);
        if (parsed) return parsed;
      }
      if (Array.isArray(profile?.courses) && profile.courses.length > 0) {
        const fallbackSubs = profile.courses.map((c: any, idx: number) => ({
          id: c.id || c._id || `course_sub_${idx}`,
          course: { id: c.id || c._id, title: c.title },
          plan_name: c.title,
          status: "active",
          start_date: c.created_at || new Date().toISOString(),
          end_date: "",
          remaining_days: 30,
          total_days: 30,
        }));
        return this.parseSubscriptionList(fallbackSubs);
      }
    } catch (err) {
      console.warn("Profile fallback in subscriptionService error:", err);
    }

    return null;
  }

  public async activateCode(code: string): Promise<{ success: boolean; message: string }> {
    if (!code.trim()) {
      return { success: false, message: "يرجى إدخال رمز الكود أولاً." };
    }
    try {
      await redeemSubscriptionCode(code);
      return { success: true, message: "تم تفعيل الاشتراك بنجاح!" };
    } catch (err: any) {
      if (err?.response?.data?.detail) {
        return { success: false, message: err.response.data.detail };
      }
      if (err?.message) {
        return { success: false, message: err.message };
      }
      return { success: false, message: "كود الاشتراك غير صالح أو منتهي الصلاحية." };
    }
  }
}

export const subscriptionService = SubscriptionService.getInstance();
