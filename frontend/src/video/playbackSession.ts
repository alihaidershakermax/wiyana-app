import {
  createPlaybackSession,
  sendHeartbeat,
  PlaybackSessionResponse,
} from "@/api/securityApi";
import { playbackTokenManager } from "@/video/playbackToken";
import { Platform } from "@/security/platformCapabilities";

export class PlaybackSessionClient {
  private sessionId: string | null = null;
  private videoId: string | null = null;
  private platform: Platform | null = null;
  private heartbeatInterval: ReturnType<typeof setInterval> | null = null;
  private heartbeatOnError: ((err: any) => void) | null = null;

  async startSession(videoId: string, courseId?: string, platform?: Platform): Promise<PlaybackSessionResponse> {
    const res = await createPlaybackSession(videoId, courseId, platform);
    this.sessionId = res.session.sessionId;
    this.videoId = videoId;
    this.platform = platform ?? null;

    playbackTokenManager.setToken(res.token, res.session.sessionId, res.session.expiresAt);
    return res;
  }

  startHeartbeat(
    getCurrentTime: () => number,
    getIsPlaying: () => boolean,
    getRiskScore?: () => number,
    onError?: (err: any) => void
  ): void {
    this.stopHeartbeat();
    this.heartbeatOnError = onError ?? null;
    this.heartbeatInterval = setInterval(async () => {
      if (!this.sessionId || !this.videoId) return;
      try {
        await sendHeartbeat(
          this.sessionId,
          this.videoId,
          getCurrentTime(),
          getIsPlaying(),
          getRiskScore?.() ?? 0,
          this.platform ?? undefined
        );
      } catch (err) {
        this.heartbeatOnError?.(err);
      }
    }, 10000); // 10s heartbeat
  }

  stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  revoke(): void {
    this.stopHeartbeat();
    playbackTokenManager.clearToken();
  }

  getSessionId(): string | null {
    return this.sessionId;
  }
}