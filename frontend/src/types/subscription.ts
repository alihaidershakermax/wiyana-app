export type SubscriptionStatus = "active" | "expired" | "none";

export interface CourseSubscriptionItem {
  id: string;
  courseTitle: string;
  startDate: string;
  endDate: string;
  remainingDays: number;
  totalDays: number;
  status: SubscriptionStatus;
}

export interface SubscriptionInfo {
  id?: string;
  status: SubscriptionStatus;
  planName?: string;
  startDate?: string;
  endDate?: string;
  remainingDays?: number;
  totalDays?: number;
  subscriptionsList?: CourseSubscriptionItem[];
}
