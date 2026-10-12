import { onBeforeUnmount, onMounted, ref } from "vue";
import { ClientRiskEngine, RiskLevel } from "@/security/riskEngine";
import { reportSecurityTelemetry } from "@/api/securityApi";
import { securityEventBus, SecurityEventType } from "@/security/securityEvents";
import {
  createCaptureDetector,
  detectPlatform,
  detectCapabilities,
  SecuritySignal,
  CaptureDetector,
} from "@/security/captureDetection";
import { Platform, PlatformCapabilities } from "@/security/platformCapabilities";
import { getPlatformPolicy } from "@/security/platformPolicy";

export interface UseCaptureDetectionOptions {
  videoId: string;
  sessionId: string;
  platform?: Platform;
  onRiskChange?: (score: number, level: RiskLevel) => void;
  onHighRisk?: (score: number) => void;
  onCritical?: (score: number) => void;
  onTampered?: () => void;
  onSignal?: (signal: SecuritySignal) => void;
}

const TELEMETRY_COOLDOWN: Record<string, number> = {};

function throttleTelemetry(eventType: string, ms = 15000): boolean {
  const now = Date.now();
  if (now - (TELEMETRY_COOLDOWN[eventType] || 0) < ms) return false;
  TELEMETRY_COOLDOWN[eventType] = now;
  return true;
}

export function useCaptureDetection(opts: UseCaptureDetectionOptions) {
  const platform = ref<Platform>(opts.platform ?? detectPlatform());
  const capabilities = ref<PlatformCapabilities>(detectCapabilities(platform.value));
  const policy = getPlatformPolicy(platform.value);

  const riskEngine = new ClientRiskEngine(policy.weights);
  let detector: CaptureDetector | null = null;

  const currentRiskScore = ref(0);
  const currentRiskLevel = ref<RiskLevel>("NORMAL");
  let lastEmitLevel: RiskLevel = "NORMAL";

  const emitSecurityEvent = (
    eventType: SecurityEventType,
    confidence: number,
    metadata: Record<string, any> = {}
  ) => {
    securityEventBus.emit({
      eventId: crypto.randomUUID(),
      userId: "",
      videoId: opts.videoId,
      sessionId: opts.sessionId,
      eventType,
      riskScore: currentRiskScore.value,
      confidence,
      timestamp: Date.now(),
      platform: platform.value,
      metadata,
    });
  };

  const updateRisk = (eventType: string, confidence = 1.0, metadata: Record<string, any> = {}) => {
    const score = riskEngine.addSignal(eventType, confidence);
    const level = riskEngine.getLevel();

    currentRiskScore.value = score;
    currentRiskLevel.value = level;

    opts.onRiskChange?.(score, level);

    // Telemetry to backend (throttled per event type to limit load).
    if (throttleTelemetry(eventType)) {
      void reportSecurityTelemetry(eventType, opts.videoId, opts.sessionId, confidence, {
        ...metadata,
        platform: platform.value,
      });
    }

    emitSecurityEvent(eventType as SecurityEventType, confidence, metadata);

    if (level !== lastEmitLevel) {
      if (level === "SUSPICIOUS") {
        emitSecurityEvent("SUSPICIOUS_ACTIVITY", Math.min(1, confidence + 0.1));
      } else if (level === "HIGH_RISK") {
        emitSecurityEvent("HIGH_RISK_ACTIVITY", Math.min(1, confidence + 0.15));
      }
      if (level === "HIGH_RISK") {
        opts.onHighRisk?.(score);
      }
      if (level === "CRITICAL") {
        opts.onCritical?.(score);
      }
      lastEmitLevel = level;
    }
  };

  const handleSignal = (signal: SecuritySignal) => {
    opts.onSignal?.(signal);
    updateRisk(signal.type, signal.confidence, {
      ...signal.metadata,
      at: signal.timestamp,
    });
  };

  const handleTamperDetected = () => {
    updateRisk("WATERMARK_TAMPER", 1.0, { tamper: true, source: "dynamic-watermark" });
    opts.onTampered?.();
  };

  const reportPlayerTamper = () => {
    updateRisk("PLAYER_TAMPER", 1.0, { source: "player-integrity" });
  };

  const reportTokenAnomaly = () => {
    updateRisk("TOKEN_ANOMALY", 1.0, { source: "playback-token" });
  };

  onMounted(() => {
    detector = createCaptureDetector(platform.value, { onSignal: handleSignal, capabilities: capabilities.value });
    detector.initialize();
  });

  onBeforeUnmount(() => {
    detector?.destroy();
    detector = null;
    riskEngine.reset();
  });

  return {
    platform: platform.value,
    capabilities: capabilities.value,
    detector,
    currentRiskScore,
    currentRiskLevel,
    updateRisk,
    handleTamperDetected,
    reportPlayerTamper,
    reportTokenAnomaly,
    resetRisk: () => {
      riskEngine.reset();
      currentRiskScore.value = 0;
      currentRiskLevel.value = "NORMAL";
      lastEmitLevel = "NORMAL";
    },
  };
}