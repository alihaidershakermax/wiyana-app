import { fetchCourseLessons, fetchLessonById, LessonApiResponse } from "@/api/lessons";
import { fetchCourseById } from "@/api/courses";
import { fetchVideoById, VideoApiResponse } from "@/api/videos";
import { fetchMySubscriptions } from "@/api/subscription";
import { apiClient } from "@/api/client";
import {
  Lesson,
  LessonAttachment,
  LessonAttachmentType,
  LessonInfoCourse,
  LessonPlayerBundle,
  LessonVideo,
} from "@/types/lesson";

const VIDEO_STATUS_READY = "READY";

function resolveVideoUrl(video: any, rawLesson?: any): LessonVideo | null {
  if (video) {
    if (video.youtube_video_id) {
      return {
        id: video.id || video.youtube_video_id,
        title: video.title || "فيديو الدرس",
        kind: "youtube",
        url: `https://www.youtube.com/watch?v=${video.youtube_video_id}`,
        durationSeconds: video.duration_seconds ?? undefined,
      };
    }

    if (video.youtube_url) {
      return {
        id: video.id || "yt_vid",
        title: video.title || "فيديو الدرس",
        kind: "youtube",
        url: video.youtube_url,
        durationSeconds: video.duration_seconds ?? undefined,
      };
    }

    if (video.file_path) {
      const isHls = /\.m3u8(\?.*)?$/i.test(video.file_path);
      const url = /^https?:\/\//i.test(video.file_path)
        ? video.file_path
        : `${apiClient.defaults.baseURL}${video.file_path.startsWith("/") ? "" : "/"}${video.file_path}`;
      return {
        id: video.id || "file_vid",
        title: video.title,
        kind: isHls ? "hls" : "file",
        url,
        durationSeconds: video.duration_seconds ?? undefined,
      };
    }
  }

  if (rawLesson) {
    if (rawLesson.youtube_video_id) {
      return {
        id: rawLesson.id,
        title: rawLesson.title,
        kind: "youtube",
        url: `https://www.youtube.com/watch?v=${rawLesson.youtube_video_id}`,
        durationSeconds: rawLesson.duration,
      };
    }

    if (rawLesson.youtube_url) {
      return {
        id: rawLesson.id,
        title: rawLesson.title,
        kind: "youtube",
        url: rawLesson.youtube_url,
        durationSeconds: rawLesson.duration,
      };
    }

    if (rawLesson.video_id) {
      return {
        id: rawLesson.video_id || rawLesson.id || rawLesson._id,
        title: rawLesson.title || "فيديو الدرس",
        kind: "file",
        url: "",
        durationSeconds: rawLesson.duration,
      };
    }
  }

  return null;
}

function resolveAttachments(attachmentIds: string[]): LessonAttachment[] {
  const result: LessonAttachment[] = [];

  attachmentIds.forEach((ref, index) => {
    const trimmed = String(ref || "").trim();
    if (!trimmed) return;

    const isUrl = /^https?:\/\//i.test(trimmed);
    const lastSegment = trimmed.split(/[/\\]/).filter(Boolean).pop() || trimmed;
    const fileName = lastSegment.split("?")[0];
    const ext = fileName.split(".").pop()?.toLowerCase() || "";
    const isImage = ["png", "jpg", "jpeg", "webp", "gif"].includes(ext);

    let type: LessonAttachmentType;
    let url: string;

    if (isUrl) {
      type = ext === "pdf" ? "pdf" : isImage ? "image" : "link";
      url = trimmed;
    } else {
      type = ext === "pdf" ? "pdf" : isImage ? "image" : "file";
      url = `${apiClient.defaults.baseURL}/uploads/${encodeURIComponent(trimmed)}`;
    }

    result.push({
      id: `att_${index}_${fileName}`,
      title: decodeURIComponent(fileName) || trimmed,
      type,
      url,
    });
  });

  return result;
}

function mapLesson(raw: any, video: LessonVideo | null): Lesson {
  const resolvedVideo = video || resolveVideoUrl(null, raw);

  return {
    id: String(raw.id || raw._id),
    courseId: String(raw.course_id || ""),
    title: raw.title || "محاضرة تعليمية",
    description: raw.description ?? undefined,
    order: Number(raw.order || 1),
    isFreePreview: Boolean(raw.is_free || raw.is_free_preview),
    isLocked: Boolean(raw.is_locked || raw.isLocked || false),
    accessAllowed: raw.access_allowed !== undefined ? Boolean(raw.access_allowed) : true,
    durationSeconds: raw.duration || resolvedVideo?.durationSeconds || undefined,
    createdAt: raw.created_at || new Date().toISOString(),
    video: resolvedVideo,
    attachments: resolveAttachments(raw.attachments || []),
    file_path: raw.file_path || undefined,
  };
}

export class LessonService {
  private static instance: LessonService;

  private constructor() {}

  public static getInstance(): LessonService {
    if (!LessonService.instance) {
      LessonService.instance = new LessonService();
    }
    return LessonService.instance;
  }

  public async getPlayerBundle(lessonId: string): Promise<LessonPlayerBundle> {
    const lessonRaw = await fetchLessonById(lessonId);
    const courseId = String(lessonRaw.course_id);

    const [lessonsRaw, courseRaw, videoData] = await Promise.all([
      fetchCourseLessons(courseId).catch(() => []),
      fetchCourseById(courseId).catch(() => ({})),
      lessonRaw.video_id
        ? fetchVideoById(lessonRaw.video_id).catch(() => null)
        : Promise.resolve(null),
    ]);

    const video = resolveVideoUrl(videoData, lessonRaw);

    const course: LessonInfoCourse = {
      id: String(courseRaw.id || courseId),
      title: courseRaw.title || "المساق التعليمي",
      subjectName: courseRaw.subject_name || courseRaw.subjectName || "المادة",
      instructorName: courseRaw.teacher_name || "أستاذ المساق",
      thumbnailUrl: courseRaw.thumbnail_url || courseRaw.coverImageUrl || undefined,
    };

    const isEnrolled = Boolean(courseRaw.is_enrolled || false);
    const rawLessonsList = Array.isArray(lessonsRaw) ? lessonsRaw : (lessonsRaw as any)?.data ?? [];
    const lessons = rawLessonsList.map((l: any) => mapLesson(l, null));

    return {
      lesson: mapLesson(lessonRaw, video),
      course,
      lessons,
      quiz: null,
      isEnrolled,
    };
  }
}

export const lessonService = LessonService.getInstance();
