import { useQuery } from "@tanstack/vue-query";
import { courseService } from "@/services/courseService";
import { Course } from "@/types/course";

export function useCoursesData() {
  const coursesQuery = useQuery<Course[], Error>({
    queryKey: ["coursesList"],
    queryFn: () => courseService.getCourses(),
    staleTime: 1000 * 60 * 5,
  });

  const metadataQuery = useQuery({
    queryKey: ["coursesFilterMetadata"],
    queryFn: () => courseService.getFilterMetadata(),
    staleTime: 1000 * 60 * 15,
  });

  return {
    courses: coursesQuery.data,
    isLoading: coursesQuery.isLoading,
    isError: coursesQuery.isError,
    refetch: coursesQuery.refetch,
    metadata: metadataQuery.data,
  };
}
