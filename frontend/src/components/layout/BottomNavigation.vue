<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Home, BookOpen, Bell, User } from "@lucide/vue";
import { useQuery } from "@tanstack/vue-query";
import { api } from "@/api";

const route = useRoute();
const router = useRouter();

const { data: unreadCount } = useQuery({
  queryKey: ["unreadNotificationsCount"],
  queryFn: api.notifications.fetchUnreadNotificationsCount,
  staleTime: 1000 * 60 * 2,
});

const navItems = computed(() => [
  { name: "الرئيسية", path: "/", icon: Home, hasBadge: false },
  { name: "الإشعارات", path: "/notifications", icon: Bell, hasBadge: Boolean(unreadCount.value && unreadCount.value > 0) },
  { name: "حسابي", path: "/profile", icon: User, hasBadge: false },
]);

function isItemActive(path: string): boolean {
  if (path === "/") {
    return route.path === "/";
  }
  return route.path.startsWith(path);
}
</script>

<template>
  <nav
    v-if="!route.path.startsWith('/teacher')"
    class="glass-liquid-nav fixed bottom-0 left-0 right-0 z-50 bg-white/75 dark:bg-[#181818]/78 backdrop-blur-xl border-t border-[#E5E7EB]/70 dark:border-[#27272A] px-2 py-1.5 pb-safe transition-all shadow-[0_-1px_3px_rgba(0,0,0,0.03)]"
  >
    <div class="max-w-md mx-auto flex items-center justify-around">
      <button
        v-for="item in navItems"
        :key="item.path"
        type="button"
        @click="router.push(item.path)"
        :aria-label="item.name"
        :aria-current="isItemActive(item.path) ? 'page' : undefined"
        :class="[
          'tap-target relative flex flex-col items-center justify-center py-1 px-3 rounded-2xl transition-all duration-200 active:scale-90 focus-visible:outline-2 focus-visible:outline-orange-500 select-none cursor-pointer',
          isItemActive(item.path)
            ? 'text-orange-500 font-bold'
            : 'text-[#6B7280] dark:text-[#A1A1AA] hover:text-[#111111] dark:hover:text-[#FFFFFF]'
        ]"
      >
        <!-- Active Background Pill -->
        <span
          v-if="isItemActive(item.path)"
          class="pointer-events-none absolute inset-0 bg-orange-500/10 dark:bg-orange-500/15 rounded-2xl animate-scale-in"
        ></span>

        <!-- Icon Container -->
        <div class="relative z-10 p-1">
          <component
            :is="item.icon"
            :class="[
              'w-5 h-5 transition-transform duration-200',
              isItemActive(item.path) ? 'scale-110 stroke-[2.5]' : 'stroke-2'
            ]"
          />
          <span
            v-if="item.hasBadge"
            class="absolute top-0.5 right-0.5 w-2 h-2 rounded-full bg-orange-500 ring-2 ring-white dark:ring-[#181818]"
          ></span>
        </div>

        <!-- Label (No Emoji) -->
        <span class="relative z-10 text-[10px] tracking-tight mt-0.5 font-bold">
          {{ item.name }}
        </span>
      </button>
    </div>
  </nav>
</template>

<style scoped>
@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.85);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
  .animate-scale-in {
  animation: scaleIn 0.2s var(--motion-snap) forwards;
}
</style>
