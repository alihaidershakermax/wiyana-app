import { ref } from "vue";
import { RiskLevel } from "@/security/riskEngine";
import {
  getPlatformPolicy,
  messageForLevel,
  unavailableMessage,
  PlatformPolicy,
} from "@/security/platformPolicy";
import { detectPlatform, Platform } from "@/security/platformCapabilities";

export interface SecurityGuardState {
  isBlocked: boolean;
  warningModalOpen: boolean;
  warningLevel: RiskLevel;
  warningTitle: string;
  warningMessage: string;
  warningDetail: string;
  requireRevalidation: boolean;
  currentLevel: RiskLevel;
}

export function useSecurityGuard(platform?: Platform) {
  const resolvedPlatform = platform ?? detectPlatform();
  const policy: PlatformPolicy = getPlatformPolicy(resolvedPlatform);

  const currentLevel = ref<RiskLevel>("NORMAL");
  const warningLevel = ref<RiskLevel>("NORMAL");
  const warningTitle = ref("");
  const warningMessage = ref("");
  const warningDetail = ref("");
  const warningModalOpen = ref(false);
  const isBlocked = ref(false);
  const requireRevalidation = ref(false);
  const detectionUnavailable = ref(false);

  const openWarning = (
    level: RiskLevel,
    title: string,
    message: string,
    detail = "",
    revalidate = false
  ) => {
    warningLevel.value = level;
    warningTitle.value = title;
    warningMessage.value = message;
    warningDetail.value = detail;
    requireRevalidation.value = revalidate;
    warningModalOpen.value = true;
    if (level === "HIGH_RISK" || level === "CRITICAL") {
      isBlocked.value = true;
    }
  };

  /** Advisory (non-blocking) notice for SUSPICIOUS activity. */
  const showSuspiciousNotice = (msg?: string) => {
    currentLevel.value = "SUSPICIOUS";
    openWarning("SUSPICIOUS", "تنبيه أمني", msg || policy.messages.suspicious);
  };

  /** Pause-and-warn for HIGH_RISK. */
  const triggerHighRiskLockout = (msg?: string) => {
    currentLevel.value = "HIGH_RISK";
    openWarning(
      "HIGH_RISK",
      "تم إيقاف المحاضرة مؤقتًا",
      msg || messageForLevel(policy, "HIGH_RISK", false),
      messageForLevel(policy, "HIGH_RISK", true),
      false
    );
  };

  /** Revoke-and-revalidate for CRITICAL. */
  const triggerCriticalLockout = (msg?: string) => {
    currentLevel.value = "CRITICAL";
    openWarning(
      "CRITICAL",
      "تم إنهاء جلسة العرض لحماية المحتوى",
      msg || messageForLevel(policy, "CRITICAL", false),
      messageForLevel(policy, "CRITICAL", true),
      true
    );
  };

  /** Non-accusatory notice when the platform cannot verify capture. */
  const showDetectionUnavailable = () => {
    detectionUnavailable.value = true;
    openWarning("SUSPICIOUS", "تنبيه أمني", unavailableMessage(policy), "", false);
  };

  const dismissWarning = () => {
    warningModalOpen.value = false;
    isBlocked.value = false;
    requireRevalidation.value = false;
  };

  return {
    policy,
    platform: resolvedPlatform,
    currentLevel,
    warningLevel,
    warningTitle,
    warningMessage,
    warningDetail,
    warningModalOpen,
    isBlocked,
    requireRevalidation,
    detectionUnavailable,
    showSuspiciousNotice,
    triggerHighRiskLockout,
    triggerCriticalLockout,
    showDetectionUnavailable,
    dismissWarning,
  };
}