import { fetchLessonQuiz } from "@/api/quizzes";
import { LessonQuiz } from "@/types/lesson";

export class QuizService {
  private static instance: QuizService;

  private constructor() {}

  public static getInstance(): QuizService {
    if (!QuizService.instance) {
      QuizService.instance = new QuizService();
    }
    return QuizService.instance;
  }

  public async getLessonQuiz(lessonId: string): Promise<LessonQuiz | null> {
    const quiz = await fetchLessonQuiz(lessonId);
    if (!quiz) return null;

    return {
      id: quiz.id,
      title: quiz.title,
      questionsCount: quiz.questions_count,
      durationMinutes: quiz.duration_minutes,
    };
  }
}

export const quizService = QuizService.getInstance();
