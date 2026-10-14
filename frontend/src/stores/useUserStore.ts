import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { UserSettings, UserProfile } from "@/types/user";
import { loginWithTelegram, loginWithToken } from "@/api/auth";
import { fetchMyProfile } from "@/api/user";
import { getAuthToken, setAuthToken, removeAuthToken } from "@/api/client";

export const useUserStore = defineStore("user", () => {
  // Authentication & session state
  const currentUser = ref<any>(null);
  const token = ref<string | null>(getAuthToken());
  const isLoading = ref<boolean>(false);
  const authError = ref<string | null>(null);

  const isAuthenticated = computed(() => Boolean(token.value || currentUser.value));

  // Modal states
  const isLogoutModalOpen = ref<boolean>(false);
  const isActivateCodeModalOpen = ref<boolean>(false);
  const isSubscriptionDetailsModalOpen = ref<boolean>(false);
  const isSavedCoursesModalOpen = ref<boolean>(false);
  const isWatchHistoryModalOpen = ref<boolean>(false);
  const isOfflineDownloadsModalOpen = ref<boolean>(false);
  const isSupportModalOpen = ref<boolean>(false);
  const isFaqModalOpen = ref<boolean>(false);
  const isPrivacyModalOpen = ref<boolean>(false);
  const isTermsModalOpen = ref<boolean>(false);
  const isAboutModalOpen = ref<boolean>(false);
  const isActiveDevicesModalOpen = ref<boolean>(false);
  const isPinLockModalOpen = ref<boolean>(false);
  const isEditProfileModalOpen = ref<boolean>(false);

  // Toast notification state
  const isToastVisible = ref<boolean>(false);
  const toastMessage = ref<string>("");
  const toastType = ref<"success" | "info" | "error">("success");
  let toastTimer: any = null;

  function showToast(message: string, type: "success" | "info" | "error" = "success") {
    toastMessage.value = message;
    toastType.value = type;
    isToastVisible.value = true;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
      isToastVisible.value = false;
    }, 3000);
  }

  // Load saved settings or defaults (White theme default)
  const savedDarkMode = false;
  const savedLang = (localStorage.getItem("wy_lang") as "ar" | "en") || "ar";
  const savedQuality = (localStorage.getItem("wy_video_quality") as UserSettings["videoQuality"]) || "1080p";
  const savedAutoPlay = localStorage.getItem("wy_auto_play") !== "false";
  const savedNotifications = localStorage.getItem("wy_notifications") !== "false";

  const settings = ref<UserSettings>({
    darkMode: savedDarkMode,
    language: savedLang,
    videoQuality: savedQuality,
    autoPlayVideo: savedAutoPlay,
    pushNotifications: savedNotifications,
  });

  // Custom profile fallback
  const customProfile = ref<Partial<UserProfile>>({
    name: localStorage.getItem("wy_user_name") || "طالب ويانا",
    telegramUsername: localStorage.getItem("wy_telegram_user") || "",
    avatarUrl: localStorage.getItem("wy_avatar_url") || "",
  });

  async function initAuth(): Promise<void> {
    isLoading.value = true;
    authError.value = null;
    try {
      // 1. If running inside Telegram Mini App
      if (typeof window !== "undefined" && (window as any).Telegram?.WebApp) {
        const tg = (window as any).Telegram.WebApp;
        tg.ready();
        tg.expand();

        // Poll for initData (Telegram may inject it asynchronously, up to 8s)
        let tgInitData = tg.initData || "";
        if (!tgInitData) {
          const deadline = Date.now() + 8000;
          while (!tgInitData && Date.now() < deadline) {
            await new Promise((r) => setTimeout(r, 500));
            tgInitData = tg.initData || "";
          }
        }
        if (tgInitData) {
          const authRes = await loginWithTelegram(tgInitData);
          if (authRes?.access_token) {
            setAuthToken(authRes.access_token);
            token.value = authRes.access_token;
          }
          if (authRes?.user) {
            currentUser.value = authRes.user;
          }
        } else {
          // Fallback: one-time token from URL (?auth=...) when initData unavailable
          // NOTE: also checked outside Telegram block below for direct links
        }
      }

      // Fallback: one-time token from URL (?auth=...) — works with or without Telegram
      // (for direct login links opened in any browser)
      if (!token.value) {
        const urlParams = new URLSearchParams(window.location.search);
        const authToken = urlParams.get("auth");
        if (authToken) {
          try {
            const authRes = await loginWithToken(authToken);
            if (authRes?.access_token) {
              setAuthToken(authRes.access_token);
              token.value = authRes.access_token;
            }
            if (authRes?.user) {
              currentUser.value = authRes.user;
            }
            // Clean the token from URL
            urlParams.delete("auth");
            const newUrl = window.location.pathname + (urlParams.toString() ? "?" + urlParams.toString() : "");
            window.history.replaceState({}, "", newUrl);
          } catch {
            // token login failed, continue without auth
          }
        }
      }

      // 2. Fetch fresh profile from /api/v1/me
      const profile = await fetchMyProfile().catch(() => null);
      if (profile) {
        currentUser.value = profile;
        if (profile.name) customProfile.value.name = profile.name;
        if (profile.username) customProfile.value.telegramUsername = profile.username;
        if (profile.avatar || profile.photo_url) customProfile.value.avatarUrl = profile.avatar || profile.photo_url;
      }
    } catch (err: any) {
      authError.value = err?.message || "Failed to initialize authentication";
      if (import.meta.env.DEV) console.warn("[Auth Init]", err);
    } finally {
      isLoading.value = false;
    }
  }

  function logout() {
    removeAuthToken();
    token.value = null;
    currentUser.value = null;
    showToast("تم تسجيل الخروج بنجاح", "info");
  }

  function updateCustomProfile(updates: Partial<UserProfile>) {
    customProfile.value = { ...customProfile.value, ...updates };
    if (updates.name) localStorage.setItem("wy_user_name", updates.name);
    if (updates.telegramUsername) localStorage.setItem("wy_telegram_user", updates.telegramUsername);
    if (updates.avatarUrl) localStorage.setItem("wy_avatar_url", updates.avatarUrl);
    showToast("تم حفظ بيانات الحساب بنجاح");
  }

  // Settings Actions
  function toggleDarkMode(value?: boolean) {
    settings.value.darkMode = value !== undefined ? value : !settings.value.darkMode;
    localStorage.setItem("wy_dark_mode", String(settings.value.darkMode));
    if (settings.value.darkMode) {
      document.documentElement.classList.add("dark");
      showToast("تم تفعيل الوضع الليلي");
    } else {
      document.documentElement.classList.remove("dark");
      showToast("تم تفعيل الوضع النهاري الأبيض");
    }
  }

  function setLanguage(lang: "ar" | "en") {
    settings.value.language = lang;
    localStorage.setItem("wy_lang", lang);
    showToast(`تم تغيير لغة التطبيق إلى ${lang === "ar" ? "العربية" : "English"}`);
  }

  function setVideoQuality(quality: UserSettings["videoQuality"]) {
    settings.value.videoQuality = quality;
    localStorage.setItem("wy_video_quality", quality);
    showToast(`تم ضبط جودة الفيديو الافتراضية إلى ${quality}`);
  }

  function toggleAutoPlay(value?: boolean) {
    settings.value.autoPlayVideo = value !== undefined ? value : !settings.value.autoPlayVideo;
    localStorage.setItem("wy_auto_play", String(settings.value.autoPlayVideo));
    showToast(settings.value.autoPlayVideo ? "تم تفعيل التشغيل التلقائي" : "تم إيقاف التشغيل التلقائي");
  }

  function toggleNotifications(value?: boolean) {
    settings.value.pushNotifications = value !== undefined ? value : !settings.value.pushNotifications;
    localStorage.setItem("wy_notifications", String(settings.value.pushNotifications));
    showToast(settings.value.pushNotifications ? "تم تفعيل استقبال الإشعارات" : "تم كتم الإشعارات");
  }

  function clearAppCache() {
    showToast("تم مسح الذاكرة المؤقتة بنجاح");
  }

  // Modal toggle methods
  function toggleLogoutModal(open?: boolean) {
    isLogoutModalOpen.value = open !== undefined ? open : !isLogoutModalOpen.value;
  }

  function toggleActivateCodeModal(open?: boolean) {
    isActivateCodeModalOpen.value = open !== undefined ? open : !isActivateCodeModalOpen.value;
  }

  function toggleSubscriptionDetailsModal(open?: boolean) {
    isSubscriptionDetailsModalOpen.value = open !== undefined ? open : !isSubscriptionDetailsModalOpen.value;
  }

  function toggleSavedCoursesModal(open?: boolean) {
    isSavedCoursesModalOpen.value = open !== undefined ? open : !isSavedCoursesModalOpen.value;
  }

  function toggleWatchHistoryModal(open?: boolean) {
    isWatchHistoryModalOpen.value = open !== undefined ? open : !isWatchHistoryModalOpen.value;
  }

  function toggleOfflineDownloadsModal(open?: boolean) {
    isOfflineDownloadsModalOpen.value = open !== undefined ? open : !isOfflineDownloadsModalOpen.value;
  }

  function toggleSupportModal(open?: boolean) {
    isSupportModalOpen.value = open !== undefined ? open : !isSupportModalOpen.value;
  }

  function toggleFaqModal(open?: boolean) {
    isFaqModalOpen.value = open !== undefined ? open : !isFaqModalOpen.value;
  }

  function togglePrivacyModal(open?: boolean) {
    isPrivacyModalOpen.value = open !== undefined ? open : !isPrivacyModalOpen.value;
  }

  function toggleTermsModal(open?: boolean) {
    isTermsModalOpen.value = open !== undefined ? open : !isTermsModalOpen.value;
  }

  function toggleAboutModal(open?: boolean) {
    isAboutModalOpen.value = open !== undefined ? open : !isAboutModalOpen.value;
  }

  function toggleActiveDevicesModal(open?: boolean) {
    isActiveDevicesModalOpen.value = open !== undefined ? open : !isActiveDevicesModalOpen.value;
  }

  function togglePinLockModal(open?: boolean) {
    isPinLockModalOpen.value = open !== undefined ? open : !isPinLockModalOpen.value;
  }

  function toggleEditProfileModal(open?: boolean) {
    isEditProfileModalOpen.value = open !== undefined ? open : !isEditProfileModalOpen.value;
  }

  return {
    // Auth State
    currentUser,
    token,
    isLoading,
    authError,
    isAuthenticated,
    initAuth,
    logout,

    // Settings & Customization
    settings,
    customProfile,
    isToastVisible,
    toastMessage,
    toastType,

    // Modal state refs
    isLogoutModalOpen,
    isActivateCodeModalOpen,
    isSubscriptionDetailsModalOpen,
    isSavedCoursesModalOpen,
    isWatchHistoryModalOpen,
    isOfflineDownloadsModalOpen,
    isSupportModalOpen,
    isFaqModalOpen,
    isPrivacyModalOpen,
    isTermsModalOpen,
    isAboutModalOpen,
    isActiveDevicesModalOpen,
    isPinLockModalOpen,
    isEditProfileModalOpen,

    // Actions
    showToast,
    updateCustomProfile,
    toggleDarkMode,
    setLanguage,
    setVideoQuality,
    toggleAutoPlay,
    toggleNotifications,
    clearAppCache,

    // Modal toggle actions
    toggleLogoutModal,
    toggleActivateCodeModal,
    toggleSubscriptionDetailsModal,
    toggleSavedCoursesModal,
    toggleWatchHistoryModal,
    toggleOfflineDownloadsModal,
    toggleSupportModal,
    toggleFaqModal,
    togglePrivacyModal,
    toggleTermsModal,
    toggleAboutModal,
    toggleActiveDevicesModal,
    togglePinLockModal,
    toggleEditProfileModal,
  };
});
