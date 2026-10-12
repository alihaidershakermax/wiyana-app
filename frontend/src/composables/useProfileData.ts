import { computed } from "vue";
import { useQuery } from "@tanstack/vue-query";
import { profileService, ProfileDataResponse } from "@/services/profileService";
import { subscriptionService } from "@/services/subscriptionService";
import { SubscriptionInfo } from "@/types/subscription";

export function useProfileData() {
  const profileQuery = useQuery<ProfileDataResponse | null, Error>({
    queryKey: ["userProfileData"],
    queryFn: () => profileService.getProfileData(),
    staleTime: 1000 * 60 * 3,
  });

  const subscriptionQuery = useQuery<SubscriptionInfo | null, Error>({
    queryKey: ["userSubscriptionInfo"],
    queryFn: () => subscriptionService.getSubscriptionInfo(),
    staleTime: 1000 * 60 * 3,
  });

  const profile = computed(() => profileQuery.data.value ?? null);

  const subscription = computed<SubscriptionInfo | null>(() => {
    // 1. Check direct subscription query
    if (subscriptionQuery.data.value) {
      return subscriptionQuery.data.value;
    }

    // 2. Check profile query subscriptions array
    const profileData = profileQuery.data.value;
    if (profileData && Array.isArray(profileData.subscriptions) && profileData.subscriptions.length > 0) {
      const parsed = subscriptionService.parseSubscriptionList(profileData.subscriptions);
      if (parsed) return parsed;
    }

    // 3. Check profile query enrolled courses
    if (profileData && Array.isArray(profileData.courses) && profileData.courses.length > 0) {
      const courseSubs = profileData.courses.map((c: any, i: number) => ({
        id: String(c.id || c._id || `enr_${i}`),
        plan_name: c.title || "مساق دراسي",
        course: { id: String(c.id || c._id || ""), title: c.title || "" },
        status: "active",
        start_date: c.created_at || new Date().toISOString(),
        end_date: "",
        remaining_days: 30,
        total_days: 30,
      }));
      return subscriptionService.parseSubscriptionList(courseSubs);
    }

    return null;
  });

  return {
    profile,
    subscription,
    isLoading: computed(() => profileQuery.isLoading.value),
    isError: computed(() => profileQuery.isError.value),
    isProfileLoading: computed(() => profileQuery.isLoading.value),
    isProfileError: computed(() => profileQuery.isError.value),
    isSubscriptionLoading: computed(() => subscriptionQuery.isLoading.value && profileQuery.isLoading.value),
    isSubscriptionError: computed(() => subscriptionQuery.isError.value && profileQuery.isError.value),
    refetchSubscription: () => {
      subscriptionQuery.refetch();
      profileQuery.refetch();
    },
    refetchProfile: () => profileQuery.refetch(),
    refetch: () => {
      profileQuery.refetch();
      subscriptionQuery.refetch();
    },
  };
}