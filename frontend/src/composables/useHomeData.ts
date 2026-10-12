import { useQuery } from "@tanstack/vue-query";
import { homeService } from "@/services/homeService";
import { HomeDashboardData } from "@/types/home";

export function useHomeData() {
  return useQuery<HomeDashboardData, Error>({
    queryKey: ["homeDashboardData"],
    queryFn: () => homeService.getHomeDashboard(),
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 2,
  });
}
