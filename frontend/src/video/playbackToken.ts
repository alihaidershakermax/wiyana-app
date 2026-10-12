export interface PlaybackToken {
  tokenString: string;
  sessionId: string;
  expiresAt: number;
}

export class PlaybackTokenManager {
  private activeToken: PlaybackToken | null = null;

  setToken(tokenString: string, sessionId: string, expiresAt: number): void {
    this.activeToken = { tokenString, sessionId, expiresAt };
  }

  getToken(): string | null {
    if (!this.activeToken) return null;
    if (Date.now() / 1000 > this.activeToken.expiresAt) {
      return null;
    }
    return this.activeToken.tokenString;
  }

  clearToken(): void {
    this.activeToken = null;
  }
}

export const playbackTokenManager = new PlaybackTokenManager();
