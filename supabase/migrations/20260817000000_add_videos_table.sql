-- Protected video catalog used by the Go backend and upload pipeline.
CREATE TABLE IF NOT EXISTS public.videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_path TEXT,
    temp_path TEXT,
    file_size BIGINT DEFAULT 0,
    mime_type VARCHAR(120) DEFAULT 'video/mp4',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    uploader_id VARCHAR(100),
    course_id VARCHAR(100),
    youtube_video_id VARCHAR(255),
    youtube_url TEXT,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc'::text, now()),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT timezone('utc'::text, now())
);

CREATE INDEX IF NOT EXISTS idx_videos_status ON public.videos(status);
CREATE INDEX IF NOT EXISTS idx_videos_course_id ON public.videos(course_id);
CREATE INDEX IF NOT EXISTS idx_videos_uploader_id ON public.videos(uploader_id);

ALTER TABLE public.videos ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Service role can manage videos"
    ON public.videos
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);
