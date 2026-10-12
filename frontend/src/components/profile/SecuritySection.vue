<script setup lang="ts">
import { LogOut } from "@lucide/vue";
import { useUserStore } from "@/stores/useUserStore";

const userStore = useUserStore();

const emit = defineEmits<{
  (e: "clickLogout"): void;
}>();

function handleLogout() {
  userStore.toggleLogoutModal(true);
  emit("clickLogout");
}
</script>

<template>
  <div class="flex flex-col gap-3 w-full">
    <h3 class="text-sm font-black text-slate-800 dark:text-slate-200">
      الأمان والحساب
    </h3>

    <div class="flex flex-col bg-white dark:bg-slate-900 border border-slate-200/80 dark:border-slate-800/80 rounded-3xl overflow-hidden shadow-xs divide-y divide-slate-100 dark:divide-slate-800/60 p-1">
      <!-- Logout Button -->
      <button
        @click="handleLogout"
        type="button"
        class="w-full flex items-center justify-between p-4 hover:bg-rose-50 dark:hover:bg-rose-950/20 transition-colors text-right group cursor-pointer active:scale-99"
      >
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-rose-500/10 text-rose-500 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform">
            <LogOut class="w-5 h-5" />
          </div>
          <div class="flex flex-col">
            <span class="text-xs sm:text-sm font-bold text-rose-600 dark:text-rose-400">
              تسجيل الخروج
            </span>
            <span class="text-[11px] text-slate-400 font-medium">
              إنهاء الجلسة الحالية على هذا الجهاز
            </span>
          </div>
        </div>
      </button>
    </div>
  </div>
</template>
