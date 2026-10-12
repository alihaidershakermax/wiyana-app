<script setup lang="ts">
import { computed } from "vue";
import { Moon, Sun, Bell, Globe, User, Check } from "@lucide/vue";
import { useUserStore } from "@/stores/useUserStore";
import GroupedListSection from "./GroupedListSection.vue";
import GroupedListItem from "./GroupedListItem.vue";

const userStore = useUserStore();

const emit = defineEmits<{
  (e: "openLanguage"): void;
  (e: "openAccount"): void;
}>();

const languageLabel = computed(() => {
  return userStore.settings.language === "en" ? "English" : "العربية";
});
</script>

<template>
  <GroupedListSection title="الإعدادات">
    <!-- 1. Appearance / Dark Mode -->
    <GroupedListItem
      title="المظهر"
      description="الوضع الليلي والنهاري"
      :icon="userStore.settings.darkMode ? Moon : Sun"
      icon-color-class="text-orange-500"
      icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
      :is-clickable="true"
      :show-chevron="false"
      @click="userStore.toggleDarkMode()"
    >
      <template #trailing>
        <!-- iOS Toggle Switch -->
        <button
          type="button"
          @click.stop="userStore.toggleDarkMode()"
          :class="[
            'w-12 h-7 rounded-full p-1 transition-all duration-300 flex items-center cursor-pointer focus:outline-none shrink-0',
            userStore.settings.darkMode ? 'bg-orange-500 justify-end' : 'bg-[#E5E7EB] dark:bg-[#27272A] justify-start'
          ]"
          :aria-checked="userStore.settings.darkMode"
          role="switch"
        >
          <div class="w-5 h-5 rounded-full bg-white shadow-sm flex items-center justify-center transition-transform">
            <Check v-if="userStore.settings.darkMode" class="w-3 h-3 text-orange-500 stroke-[3]" />
          </div>
        </button>
      </template>
    </GroupedListItem>

    <!-- 2. Notifications -->
    <GroupedListItem
      title="الإشعارات"
      description="تنبيهات الدروس والنشاط"
      :icon="Bell"
      icon-color-class="text-orange-500"
      icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
      :is-clickable="true"
      :show-chevron="false"
      @click="userStore.toggleNotifications()"
    >
      <template #trailing>
        <!-- iOS Toggle Switch -->
        <button
          type="button"
          @click.stop="userStore.toggleNotifications()"
          :class="[
            'w-12 h-7 rounded-full p-1 transition-all duration-300 flex items-center cursor-pointer focus:outline-none shrink-0',
            userStore.settings.pushNotifications ? 'bg-orange-500 justify-end' : 'bg-[#E5E7EB] dark:bg-[#27272A] justify-start'
          ]"
          :aria-checked="userStore.settings.pushNotifications"
          role="switch"
        >
          <div class="w-5 h-5 rounded-full bg-white shadow-sm flex items-center justify-center transition-transform">
            <Check v-if="userStore.settings.pushNotifications" class="w-3 h-3 text-orange-500 stroke-[3]" />
          </div>
        </button>
      </template>
    </GroupedListItem>

    <!-- 3. Language -->
    <GroupedListItem
      title="اللغة"
      description="لغة واجهة التطبيق"
      :icon="Globe"
      icon-color-class="text-orange-500"
      icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
      :value-text="languageLabel"
      :show-chevron="true"
      @click="emit('openLanguage')"
    />

    <!-- 4. Account -->
    <GroupedListItem
      title="الحساب"
      description="بيانات الحساب الشخصي"
      :icon="User"
      icon-color-class="text-orange-500"
      icon-bg-class="bg-orange-500/10 dark:bg-orange-500/15"
      :show-chevron="true"
      @click="emit('openAccount')"
    />
  </GroupedListSection>
</template>
