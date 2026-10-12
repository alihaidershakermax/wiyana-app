import { UserProfile } from "./user";
import { Course } from "./course";
import { Announcement } from "./announcement";

export interface HomeDashboardData {
  user: UserProfile;
  announcement?: Announcement;
  continueLearningCourses: Course[];
  recommendedCourses: Course[];
  latestCourses: Course[];
  popularCourses: Course[];
}
