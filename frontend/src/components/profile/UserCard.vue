<script setup lang="ts">
import { computed, ref } from "vue";
import { UserProfile } from "@/types/user";
import { User, Copy, Check, Send, Calendar } from "@lucide/vue";

const props = defineProps<{
  user: UserProfile;
}>();

const emit = defineEmits<{
  (e: "copyId", id: string): void;
}>();

const isCopied = ref(false);

const initials = computed(() => {
  if (!props.user?.name) return "ط";
  const parts = props.user.name.trim().split(" ");
  if (parts.length >= 2) {
    return `${parts[0][0]}${parts[1][0]}`;
  }
  return props.user.name.slice(0, 2);
});

function handleCopy() {
  if (!props.user?.telegramId) return;
  navigator.clipboard?.writeText(props.user.telegramId);
  isCopied.value = true;
  emit("copyId", props.user.telegramId);
  setTimeout(() => {
    isCopied.value = false;
  }, 2000);
}
</script>

<template>
  <div
    class="w-full bg-white dark:bg-[#141414] border border-slate-200/80 dark:border-zinc-800/80 rounded-2xl sm:rounded-3xl p-5 sm:p-6 shadow-[0_1px_3px_rgba(0,0,0,0.04)] flex flex-col items-center text-center gap-4 transition-colors"
  >
    <!-- Large iOS Avatar -->
    <div class="relative">
      <div
        class="w-20 h-20 sm:w-24 sm:h-24 rounded-full bg-slate-100 dark:bg-zinc-800 ring-4 ring-orange-500/20 dark:ring-orange-500/30 overflow-hidden flex items-center justify-center shadow-sm"
      >
        <img
          v-if="user.avatarUrl"
          :src="user.avatarUrl"
          :alt="user.name"
          class="w-full h-full object-cover"
        />
        <div
          v-else-if="user.name"
          class="w-full h-full bg-gradient-to-br from-orange-500 to-orange-600 text-white font-black text-xl sm:text-2xl flex items-center justify-center select-none"
        >
          {{ initials }}
        </div>
        <User v-else class="w-10 h-10 text-slate-400 dark:text-zinc-500" />
      </div>
    </div>

    <!-- User Information -->
    <div class="flex flex-col items-center gap-1.5 w-full">
      <h2 class="text-lg sm:text-xl font-black text-slate-900 dark:text-zinc-100 tracking-tight">
        {{ user.name || "طالب ويانا" }}
      </h2>

      <!-- Telegram Details & Badges -->
      <div class="flex items-center gap-2 flex-wrap justify-center mt-1">
        <!-- Telegram Username -->
        <a
          v-if="user.telegramUsername"
          :href="`https://t.me/${user.telegramUsername}`"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-1.5 text-xs font-bold text-orange-600 dark:text-orange-400 bg-orange-500/10 hover:bg-orange-500/20 px-3 py-1.5 rounded-xl transition-colors cursor-pointer"
        >
          <Send class="w-3.5 h-3.5" />
          <span dir="ltr">@{{ user.telegramUsername }}</span>
        </a>

        <!-- Telegram ID with Instant Copy Feedback -->
        <button
          v-if="user.telegramId"
          type="button"
          @click="handleCopy"
          class="inline-flex items-center gap-1.5 text-xs font-bold text-slate-700 dark:text-zinc-300 bg-slate-100 dark:bg-zinc-800 hover:bg-slate-200 dark:hover:bg-zinc-700 px-3 py-1.5 rounded-xl transition-all cursor-pointer active:scale-95"
          title="انقر لنسخ المعرف"
        >
          <component :is="isCopied ? Check : Copy" :class="['w-3.5 h-3.5', isCopied ? 'text-emerald-500' : 'text-orange-500']" />
          <span dir="ltr">ID: {{ user.telegramId }}</span>
        </button>

        <!-- Joined Date (if available) -->
        <div
          v-if="user.joinedDate"
          class="inline-flex items-center gap-1.5 text-xs font-semibold text-slate-500 dark:text-zinc-400 bg-slate-50 dark:bg-zinc-800/60 px-3 py-1.5 rounded-xl border border-slate-200/50 dark:border-zinc-700/50"
        >
          <Calendar class="w-3.5 h-3.5 text-slate-400 dark:text-zinc-500" />
          <span>انضم في {{ user.joinedDate }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
