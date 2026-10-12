<script setup lang="ts">
import { ref, computed } from "vue";
import { useProfileData } from "@/composables/useProfileData";
import { useUserStore } from "@/stores/useUserStore";

// Icons (Strictly professional Lucide icons, Zero emojis)
import {
  HelpCircle,
  ShieldCheck,
  FileText,
  Info,
  Sparkles,
} from "@lucide/vue";

// UI & Layout Components
import ProfileHeader from "@/components/profile/ProfileHeader.vue";
import UserCard from "@/components/profile/UserCard.vue";
import SubscriptionCard from "@/components/profile/SubscriptionCard.vue";
import SettingsSection from "@/components/profile/SettingsSection.vue";
import GroupedListSection from "@/components/profile/GroupedListSection.vue";
import GroupedListItem from "@/components/profile/GroupedListItem.vue";
import SkeletonProfile from "@/components/profile/SkeletonProfile.vue";
import OfflineState from "@/components/ui/OfflineState.vue";

// Modals & Feedback
import ToastNotification from "@/components/profile/ToastNotification.vue";
import ActivateCodeModal from "@/components/profile/ActivateCodeModal.vue";
import SupportModal from "@/components/profile/SupportModal.vue";
import PrivacyModal from "@/components/profile/PrivacyModal.vue";
import TermsModal from "@/components/profile/TermsModal.vue";
import AboutModal from "@/components/profile/AboutModal.vue";
import LanguageModal from "@/components/profile/LanguageModal.vue";
import AccountDetailsModal from "@/components/profile/AccountDetailsModal.vue";

const userStore = useUserStore();
const {
  profile,
  subscription,
  isLoading,
  isSubscriptionLoading,
  isSubscriptionError,
  refetchSubscription,
  refetch,
} = useProfileData();

const isNotAuthenticated = computed(() => !userStore.isAuthenticated && !userStore.currentUser);

// Local modal states
const isLanguageModalOpen = ref(false);
const isAccountDetailsModalOpen = ref(false);

// Merged live profile object (Real API data only when authenticated)
const safeUser = computed(() => {
  const base = profile.value?.user || null;
  if (!base) {
    if (userStore.isAuthenticated && userStore.currentUser) {
      const u = userStore.currentUser;
      return {
        id: String(u.id || u._id || ""),
        telegramId: String(u.telegram_id || u.telegramId || ""),
        telegramUsername: u.username || u.telegram_username || "",
        name: `${u.first_name || ""} ${u.last_name || ""}`.trim() || u.username || u.name || "طالب ويانا",
        avatarUrl: u.photo_url || u.avatar || "",
        departmentName: u.department_name || "",
        gradeName: u.grade_name || "",
        joinedDate: u.created_at ? new Date(u.created_at).toISOString().split("T")[0] : "",
        unreadNotificationsCount: 0,
      };
    }
    return null;
  }

  return {
    ...base,
    ...userStore.customProfile,
  };
});

const progressPercentage = computed(() => {
  return profile.value?.stats?.overallCompletionPercentage ?? 0;
});
</script>

<template>
  <div
    class="min-h-screen bg-white dark:bg-[#0F0F0F] text-[#111111] dark:text-[#FFFFFF] flex flex-col font-sans transition-colors duration-200 pb-24 selection:bg-orange-500 selection:text-white"
  >
    <!-- 1. Sticky iOS Profile Header -->
    <ProfileHeader title="حسابي" />

    <!-- Floating Toast Notification -->
    <ToastNotification
      :show="userStore.isToastVisible"
      :message="userStore.toastMessage"
      :type="userStore.toastType"
      @close="userStore.isToastVisible = false"
    />

    <!-- Offline Connection State -->
    <OfflineState @refresh="refetch()" />

    <!-- Main Container -->
    <main class="flex-1 max-w-xl mx-auto px-4 py-4 sm:py-6 w-full flex flex-col gap-4">
      <!-- Loading Skeleton UI -->
      <SkeletonProfile v-if="isLoading && !safeUser" />

      <!-- Not Authenticated Prompt -->
      <div
        v-else-if="isNotAuthenticated && !isLoading"
        class="flex flex-col items-center justify-center gap-4 py-16 px-6 text-center animate-fade-in"
      >
        <div class="w-16 h-16 rounded-3xl bg-orange-500/10 flex items-center justify-center">
          <svg class="w-8 h-8 text-orange-500" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold text-[#111111] dark:text-white">تسجيل الدخول مطلوب</h3>
          <p class="text-sm text-[#6B7280] dark:text-[#A1A1AA] mt-2 leading-relaxed">
            يرجى فتح التطبيق عبر بوت تليجرام لعرض بيانات حسابك والاشتراكات.
          </p>
        </div>
        <a
          href="https://t.me/WiyanaAppBot"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 px-6 py-3 rounded-2xl bg-orange-500 hover:bg-[#EA580C] active:scale-95 text-white font-bold text-sm shadow-sm shadow-orange-500/20 transition-all"
        >
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.07-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .24z"/>
          </svg>
          فتح بوت تليجرام
        </a>
      </div>

      <!-- Loaded Content -->
      <div v-else class="flex flex-col gap-4 animate-fade-in">
        <!-- 2. User Profile Card (Avatar, Name, @Username, ID, Join Date) -->
        <UserCard
          v-if="safeUser"
          :user="safeUser"
          @copy-id="(id) => userStore.showToast(`تم نسخ معرف تيليجرام ID: ${id}`)"
        />

        <!-- 3. Subscription Status Card (Active / Expired / None / Error / Skeleton) -->
        <SubscriptionCard
          :subscription="subscription"
          :is-loading="isSubscriptionLoading"
          :is-error="isSubscriptionError"
          :progress-percentage="progressPercentage"
          @activate-code="userStore.toggleActivateCodeModal(true)"
          @retry="refetchSubscription()"
        />

        <!-- 4. Settings Section (Appearance, Notifications, Language, Account) -->
        <SettingsSection
          @open-language="isLanguageModalOpen = true"
          @open-account="isAccountDetailsModalOpen = true"
        />

        <!-- 5. Support Section -->
        <GroupedListSection title="الدعم">
          <GroupedListItem
            title="مركز الدعم"
            description="تحتاج مساعدة أو تواجه مشكلة؟"
            :icon="HelpCircle"
            icon-color-class="text-orange-500"
            icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
            :show-chevron="true"
            @click="userStore.toggleSupportModal(true)"
          />
        </GroupedListSection>

        <!-- 6. Legal Information Section -->
        <GroupedListSection title="المعلومات القانونية">
          <GroupedListItem
            title="سياسة الخصوصية"
            description="حماية البيانات وسرية الحساب"
            :icon="ShieldCheck"
            icon-color-class="text-orange-500"
            icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
            :show-chevron="true"
            @click="userStore.togglePrivacyModal(true)"
          />

          <GroupedListItem
            title="شروط الاستخدام"
            description="قواعد وأحكام المنصة التعليمية"
            :icon="FileText"
            icon-color-class="text-orange-500"
            icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
            :show-chevron="true"
            @click="userStore.toggleTermsModal(true)"
          />
        </GroupedListSection>

        <!-- 7. About App Section -->
        <GroupedListSection title="عن التطبيق">
          <GroupedListItem
            title="حول WEYANA"
            description="منصة تعليمية تساعد الطلاب على الوصول إلى مساقاتهم ومحاضراتهم ومتابعة تقدمهم الدراسي."
            :icon="Info"
            icon-color-class="text-orange-500"
            icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
            value-text="الإصدار 1.0.0"
            :show-chevron="true"
            @click="userStore.toggleAboutModal(true)"
          />
        </GroupedListSection>
      </div>
    </main>

    <!-- ── Interactive iOS Modals ────────────────────────────────────────── -->

    <!-- Activate Subscription Code Modal -->
    <ActivateCodeModal
      :is-open="userStore.isActivateCodeModalOpen"
      @close="userStore.toggleActivateCodeModal(false)"
      @activated="() => {
        refetch();
        userStore.showToast('تم تفعيل اشتراكك بنجاح! جميع المحاضرات مفتوحة الآن.');
      }"
    />

    <!-- Support Center Modal -->
    <SupportModal
      :is-open="userStore.isSupportModalOpen"
      @close="userStore.toggleSupportModal(false)"
      @sent="userStore.showToast('تم إرسال تذكرة الدعم بنجاح! سيتم الرد عليك.')"
    />

    <!-- Privacy Policy Modal -->
    <PrivacyModal
      :is-open="userStore.isPrivacyModalOpen"
      @close="userStore.togglePrivacyModal(false)"
    />

    <!-- Terms of Service Modal -->
    <TermsModal
      :is-open="userStore.isTermsModalOpen"
      @close="userStore.toggleTermsModal(false)"
    />

    <!-- About WEYANA Platform Modal -->
    <AboutModal
      :is-open="userStore.isAboutModalOpen"
      @close="userStore.toggleAboutModal(false)"
    />

    <!-- Language Selector Modal -->
    <LanguageModal
      :is-open="isLanguageModalOpen"
      @close="isLanguageModalOpen = false"
    />

    <!-- Account Details Modal -->
    <AccountDetailsModal
      :is-open="isAccountDetailsModalOpen"
      :user="safeUser"
      @close="isAccountDetailsModalOpen = false"
    />
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
