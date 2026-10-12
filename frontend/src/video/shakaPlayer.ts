import shaka from "shaka-player";

export interface ShakaPlayerConfig {
  videoElement: HTMLVideoElement;
  manifestUri: string;
  token?: string;
  onRateChanged?: (rate: number) => void;
  onError?: (error: any) => void;
}

export class WeyanaShakaPlayer {
  private player: shaka.Player | null = null;
  private videoElement: HTMLVideoElement;

  constructor(config: ShakaPlayerConfig) {
    this.videoElement = config.videoElement;
    shaka.polyfill.installAll();

    if (shaka.Player.isBrowserSupported()) {
      this.initPlayer(config);
    } else {
      config.onError?.(new Error("Shaka Player is not supported in this browser"));
    }
  }

  private initPlayer(config: ShakaPlayerConfig): void {
    try {
      this.player = new shaka.Player(this.videoElement);

      // Disable download controls & context menu
      (this.videoElement as any).controlsList = "nodownload noplaybackrate";
      this.videoElement.disablePictureInPicture = true;

      // Lock playback rate to 1.0
      this.videoElement.addEventListener("ratechange", () => {
        if (this.videoElement.playbackRate !== 1.0) {
          this.videoElement.playbackRate = 1.0;
          config.onRateChanged?.(this.videoElement.playbackRate);
        }
      });

      // Attach token header request filter if token provided
      if (config.token) {
        this.player.getNetworkingEngine()?.registerRequestFilter((_type, request) => {
          request.headers["Authorization"] = `Bearer ${config.token}`;
        });
      }

      this.player.addEventListener("error", (event) => {
        config.onError?.(event);
      });

      void this.player.load(config.manifestUri);
    } catch (err) {
      config.onError?.(err);
    }
  }

  play(): Promise<void> {
    return this.videoElement.play();
  }

  pause(): void {
    this.videoElement.pause();
  }

  destroy(): void {
    if (this.player) {
      void this.player.destroy();
      this.player = null;
    }
  }
}
