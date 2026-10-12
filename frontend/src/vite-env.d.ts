/// <reference types="vite/client" />

declare module "*.vue" {
  import type { DefineComponent } from "vue";
  const component: DefineComponent<{}, {}, any>;
  export default component;
}

// The published screenshot-prevention package ships a broken package.json
// (its `main`/`types` point at files that do not exist) and no usable type
// declarations are resolvable. We alias the runtime to the real ESM build
// (see vite.config.ts) and declare the module here, mirroring the library's
// own dist/types/index.d.ts.
declare module "screenshot-prevention" {
  export type DetectionMethod =
    | "keyboard"
    | "mobile"
    | "screenCapture"
    | "mediaRecording"
    | "visibilityChange"
    | "devTools";

  export interface AttemptDetails {
    count: number;
    method: DetectionMethod;
    timestamp: number;
    details?: string;
  }

  export interface ScreenshotPreventionOptions {
    blurIntensity?: string;
    warningMessage?: string;
    preventCopy?: boolean;
    preventInspect?: boolean;
    recoveryDelay?: number;
    debug?: boolean;
    onAttempt?: (details: AttemptDetails) => void;
    customStyles?: Record<string, string>;
  }

  export default class EnhancedScreenshotPrevention {
    constructor(options?: Partial<ScreenshotPreventionOptions>);
    getAttemptCount(): number;
    reset(): void;
    update(options: Partial<ScreenshotPreventionOptions>): void;
    destroy(): void;
  }
}

interface TelegramWebApp {
  initData: string;
  initDataUnsafe: any;
  platform: string;
  colorScheme: string;
  isFullscreen?: boolean;
  requestFullscreen?: () => void | Promise<unknown>;
  exitFullscreen?: () => void | Promise<unknown>;
  close?: () => void;
  onEvent?: (event: string, callback: (eventData?: any) => void) => void;
  offEvent?: (event: string, callback: (eventData?: any) => void) => void;
}

interface Window {
  Telegram?: { WebApp?: TelegramWebApp };
}
