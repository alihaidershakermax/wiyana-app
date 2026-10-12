<script setup lang="ts">
import { computed } from "vue";
import { Bell, Search, GraduationCap } from "@lucide/vue";
import { UserProfile } from "@/types/user";

const props = defineProps<{
  user?: UserProfile;
}>();

const emit = defineEmits<{
  (e: "clickSearch"): void;
  (e: "clickNotifications"): void;
  (e: "clickProfile"): void;
}>();

// Dynamic time-based greeting
const currentGreeting = computed(() => {
  const hour = new Date().getHours();
  if (hour >= 4 && hour < 12) {
    return "صباح الخير";
  }
  return "مساء الخير";
});

// Clean user display name
const displayName = computed(() => {
  if (!props.user?.name || props.user.name === "زائر المنصة" || props.user.name === "طالب ويانا") {
    return "طالبنا العزيز";
  }
  return props.user.name.split(" ")[0] || props.user.name;
});
</script>

<template>
  <header class="sticky top-0 z-40 w-full backdrop-blur-xl bg-white/90 dark:bg-[#0F0F0F]/90 border-b border-[#E5E7EB] dark:border-[#27272A] transition-colors">
    <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between gap-3">
      <!-- Logo & Dynamic Greeting with Person Name -->
      <div class="flex items-center gap-3 min-w-0">
        <div class="w-11 h-11 rounded-2xl bg-gradient-to-tr from-orange-600 to-orange-500 flex items-center justify-center text-white shadow-md shadow-orange-500/25 active:scale-95 transition-transform shrink-0">
          <GraduationCap class="w-6 h-6" />
        </div>
        <div class="flex flex-col text-start min-w-0">
          <div class="flex items-center gap-1.5 min-w-0">
            <span class="text-sm sm:text-base font-black tracking-tight text-[#111111] dark:text-[#FFFFFF] truncate">
              {{ currentGreeting }}، {{ displayName }}
            </span>
            <span class="text-sm shrink-0 inline-block animate-pulse">👋</span>
          </div>
          <span class="text-[11px] font-bold text-orange-500 dark:text-orange-400 truncate -mt-0.5">
            منصة ويانا التعليمية
          </span>
        </div>
      </div>

      <!-- Action Icons -->
      <div class="flex items-center gap-2 shrink-0">
        <!-- Search Button -->
        <button
          @click="emit('clickSearch')"
          type="button"
          aria-label="البحث"
          class="w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] hover:bg-[#E5E7EB] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA] flex items-center justify-center transition-all active:scale-90 focus:outline-none cursor-pointer"
        >
          <Search class="w-5 h-5" />
        </button>

        <!-- Notifications Button -->
        <button
          @click="emit('clickNotifications')"
          type="button"
          aria-label="الإشعارات"
          class="relative w-11 h-11 sm:w-12 sm:h-12 rounded-2xl bg-[#F8F8F8] dark:bg-[#181818] hover:bg-[#E5E7EB] dark:hover:bg-[#222222] border border-[#E5E7EB] dark:border-[#27272A] text-[#6B7280] dark:text-[#A1A1AA] flex items-center justify-center transition-all active:scale-90 focus:outline-none cursor-pointer"
        >
          <Bell class="w-5 h-5" />
          <span
            v-if="user?.unreadNotificationsCount && user.unreadNotificationsCount > 0"
            class="absolute top-2.5 right-2.5 w-2.5 h-2.5 rounded-full bg-orange-500 ring-2 ring-white dark:ring-[#181818]"
          ></span>
        </button>

        <!-- User Avatar Button -->
        <button
          @click="emit('clickProfile')"
          type="button"
          aria-label="الملف الشخصي"
          class="w-11 h-11 sm:w-12 sm:h-12 rounded-full bg-[#F8F8F8] dark:bg-[#181818] overflow-hidden ring-2 ring-orange-500/30 hover:ring-orange-500 transition-all active:scale-90 shrink-0 focus:outline-none cursor-pointer"
        >
          <img
            v-if="user?.avatarUrl"
            :src="user.avatarUrl"
            :alt="user.name"
            class="w-full h-full object-cover"
          />
          <div v-else class="w-full h-full flex items-center justify-center bg-orange-500 text-white font-bold text-sm">
            {{ user?.name?.charAt(0) || "و" }}
          </div>
        </button>
      </div>
    </div>
  </header>
</template>
