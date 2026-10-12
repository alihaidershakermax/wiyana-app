import { computed, toValue, type MaybeRefOrGetter } from "vue";
import { useQuery } from "@tanstack/vue-query";
import { lessonService } from "@/services/lessonService";
import { quizService } from "@/services/quizService";
import { LessonPlayerBundle } from "@/types/lesson";

export function useLessonPlayer(lessonId: MaybeRefOrGetter<string>) {
  const lessonIdRef = computed(() => toValue(lessonId));

  const playerQuery = useQuery<LessonPlayerBundle, Error>({
    queryKey: computed(() => ["lessonPlayer", lessonIdRef.value]),
    queryFn: () => lessonService.getPlayerBundle(lessonIdRef.value),
    staleTime: 1000 * 60 * 5,
    retry: 1,
  });

  const quizQuery = useQuery({
    queryKey: computed(() => ["lessonQuiz", lessonIdRef.value]),
    queryFn: () => quizService.getLessonQuiz(lessonIdRef.value),
    staleTime: 1000 * 60 * 5,
    retry: false,
  });

  return {
    bundle: playerQuery.data,
    isLoading: playerQuery.isLoading,
    isError: playerQuery.isError,
    isRefetching: playerQuery.isFetching,
    refetch: playerQuery.refetch,
    quiz: quizQuery.data,
  };
}
