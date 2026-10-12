import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { ContentProtectionEngine } from "@/services/contentProtection";
import type { AccessStatus, ViolationType } from "@/types/protection";

export function useContentProtection(opts: {
  lessonId?: string;
  closeOnCapture?: boolean;
  onBanned?: (result: any) => void;
  onWarning?: (result: any) => void;
  onViolation?: (type: ViolationType, detectedBy: string) => void;
  onPrivacyPause?: (detectedBy: string) => void;
}) {
  const engine = new ContentProtectionEngine({
    lessonId: opts.lessonId,
    closeOnCapture: opts.closeOnCapture,
    onBanned: opts.onBanned,
    onWarning: opts.onWarning,
    onViolation: opts.onViolation,
    onPrivacyPause: opts.onPrivacyPause,
  });

  const status = ref<AccessStatus | null>(null);
  const checked = ref(false);

  const banned = computed(() => status.value?.banned ?? false);
  const allowed = computed(() => checked.value && (status.value?.allowed ?? true));

  onMounted(async () => {
    engine.start();
    try {
      status.value = await engine.checkAccessNow();
      if (status.value?.banned) {
        opts.onBanned?.({
          success: true,
          points: status.value.points,
          status: "banned",
          banned: true,
          warning_threshold: status.value.warning_threshold,
          ban_threshold: status.value.ban_threshold,
          ban_reason: status.value.ban_reason,
          ban_until: status.value.ban_until,
          ban_days: status.value.ban_days,
        });
      }
    } catch {
      status.value = null;
    } finally {
      checked.value = true;
    }
  });

  onBeforeUnmount(() => {
    engine.stop();
  });

  return { engine, status, checked, banned, allowed };
}
