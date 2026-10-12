export interface Announcement {
  id: string;
  title: string;
  description: string;
  actionText?: string;
  actionUrl?: string;
  type: "info" | "warning" | "success" | "promotion";
  isDismissible: boolean;
}
