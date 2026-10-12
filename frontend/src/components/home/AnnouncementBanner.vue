<script setup lang="ts">
import { Megaphone, X, ArrowLeft } from "@lucide/vue";
import { Announcement } from "@/types/announcement";

defineProps<{
  announcement?: Announcement;
}>();

const emit = defineEmits<{
  (e: "dismiss"): void;
  (e: "clickAction", url?: string): void;
}>();
</script>

<template>
  <div
    v-if="announcement"
    class="relative w-full rounded-2xl bg-gradient-to-r from-orange-500/10 via-orange-500/10 to-orange-500/5 border border-orange-500/30 p-4 sm:p-5 flex items-start sm:items-center justify-between gap-4 shadow-sm"
  >
    <div class="flex items-start sm:items-center gap-3.5 flex-1">
      <div class="w-10 h-10 rounded-xl bg-orange-500 text-white flex items-center justify-center shrink-0 shadow-md shadow-orange-500/20">
        <Megaphone class="w-5 h-5" />
      </div>

      <div class="flex flex-col gap-0.5">
        <h3 class="text-sm sm:text-base font-bold text-slate-900 dark:text-slate-100">
          {{ announcement.title }}
        </h3>
        <p class="text-xs sm:text-sm text-slate-600 dark:text-slate-300">
          {{ announcement.description }}
        </p>
      </div>
    </div>

    <div class="flex items-center gap-2 shrink-0">
      <button
        v-if="announcement.actionText"
        @click="emit('clickAction', announcement.actionUrl)"
        type="button"
        class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-orange-500 hover:bg-orange-600 text-white text-xs font-bold transition-all shadow-sm shadow-orange-500/20"
      >
        <span>{{ announcement.actionText }}</span>
        <ArrowLeft class="w-3.5 h-3.5" />
      </button>

      <button
        v-if="announcement.isDismissible"
        @click="emit('dismiss')"
        type="button"
        aria-label="إغلاق الإعلان"
        class="p-2 rounded-xl text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-200/50 dark:hover:bg-slate-800 transition-colors"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
