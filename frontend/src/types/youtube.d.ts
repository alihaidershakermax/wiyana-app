export {};

declare global {
  interface Window {
    YT?: {
      Player: new (
        element: HTMLElement,
        options: Record<string, any>
      ) => {
        playVideo(): void;
        pauseVideo(): void;
        seekTo(seconds: number, allowSeekAhead: boolean): void;
        getCurrentTime(): number;
        getDuration(): number;
        setPlaybackRate(rate: number): void;
        destroy(): void;
      };
    };
    onYouTubeIframeAPIReady?: () => void;
  }
}
