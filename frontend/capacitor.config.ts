import type { CapacitorConfig } from "@capacitor/cli";

const config: CapacitorConfig = {
  appId: "com.wiyana.app",
  appName: "Wiyana",
  webDir: "dist",
  bundledWebRuntime: false,
  plugins: {
    PrivacyScreen: {
      enable: true,
      privacyMode: "black",
    },
  },
  ios: {
    contentInset: "always",
    allowsLinkPreview: false,
    scrollEnabled: true,
  },
  android: {
    allowMixedContent: true,
    captureInput: true,
    webContentsDebuggingEnabled: false,
  },
};

export default config;
