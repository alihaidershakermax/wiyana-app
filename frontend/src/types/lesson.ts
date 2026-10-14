export type VideoKind = "youtube" | "file" | "hls";

export interface LessonVideo {
  id: string;
  title?: string;
  kind: VideoKind;
  url: string;
  durationSeconds?: number;
  youtube_video_id?: string | null;
  youtube_url?: string | null;
}

export type LessonAttachmentType = "pdf" | "file" | "image" | "link" | "telegram_file";

export interface LessonAttachment {
  id: string;
  title: string;
  type: LessonAttachmentType;
  url?: string;
  sizeBytes?: number;
}

export interface Lesson {
  id: string;
  courseId: string;
  title: string;
  description?: string;
  order: number;
  isFreePreview: boolean;
  isLocked?: boolean;
  accessAllowed?: boolean;
  isCompleted?: boolean;
  durationSeconds?: number;
  createdAt: string;
  video: LessonVideo | null;
  attachments: LessonAttachment[];
  youtube_video_id?: string | null;
  youtube_url?: string | null;
  file_path?: string | null;
}

export interface LessonInfoCourse {
  id: string;
  title: string;
  subjectName: string;
  instructorName: string;
  thumbnailUrl?: string;
}

export interface LessonQuiz {
  id: string;
  title: string;
  questionsCount: number;
  durationMinutes?: number;
}

export interface LessonPlayerBundle {
  lesson: Lesson;
  course: LessonInfoCourse;
  lessons: Lesson[];
  quiz: LessonQuiz | null;
  isEnrolled: boolean;
}
